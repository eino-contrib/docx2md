package docx_parser

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ReadDocx reads and parses a .docx file, returning a Document object.
// It now handles paragraphs, tables, headers, and footers.
func ReadDocx(filePath string) (*Document, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var doc Document
	var relationships Relationships
	var stylesList *Styles
	var relsFile, stylesFile, docFile *zip.File
	var headerFiles, footerFiles []*zip.File

	for _, f := range r.File {
		switch {
		case f.Name == "word/_rels/document.xml.rels":
			if relsFile != nil {
				return nil, fmt.Errorf("duplicate file found in docx archive: %s", f.Name)
			}
			relsFile = f
		case f.Name == "word/styles.xml":
			if stylesFile != nil {
				return nil, fmt.Errorf("duplicate file found in docx archive: %s", f.Name)
			}
			stylesFile = f
		case f.Name == "word/document.xml":
			if docFile != nil {
				return nil, fmt.Errorf("duplicate file found in docx archive: %s", f.Name)
			}
			docFile = f
		case strings.HasPrefix(f.Name, "word/header") && strings.HasSuffix(f.Name, ".xml"):
			headerFiles = append(headerFiles, f)
		case strings.HasPrefix(f.Name, "word/footer") && strings.HasSuffix(f.Name, ".xml"):
			footerFiles = append(footerFiles, f)
		}
	}

	if docFile == nil {
		return nil, fmt.Errorf("word/document.xml not found in docx")
	}
	if relsFile == nil {
		return nil, fmt.Errorf("word/_rels/document.xml.rels not found in docx")
	}

	// --- Parse files in dependency order ---

	// Parse relationships
	rcRels, err := relsFile.Open()
	if err != nil {
		return nil, err
	}
	defer rcRels.Close()
	if err := xml.NewDecoder(rcRels).Decode(&relationships); err != nil {
		return nil, err
	}

	// Parse styles (if it exists)
	if stylesFile != nil {
		// ReadStyle already handles its own file opening/closing and defer.
		stylesList, err = ReadStyle(r, filePath)
		if err != nil {
			return nil, err
		}
	}

	// Helper function to parse a content file.
	parseFile := func(f *zip.File) (Body, error) {
		rc, err := f.Open()
		if err != nil {
			return Body{}, err
		}
		defer rc.Close()
		return parseContent(xml.NewDecoder(rc), r, relationships)
	}

	// Parse main document body
	doc.Body, err = parseFile(docFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document.xml: %w", err)
	}

	// Parse headers
	for _, f := range headerFiles {
		body, err := parseFile(f)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", f.Name, err)
		}
		doc.Headers = append(doc.Headers, body)
	}

	// Parse footers
	for _, f := range footerFiles {
		body, err := parseFile(f)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", f.Name, err)
		}
		doc.Footers = append(doc.Footers, body)
	}

	// --- Apply styles to all parsed content ---
	stylizedBody(&doc.Body, stylesList)
	for i := range doc.Headers {
		stylizedBody(&doc.Headers[i], stylesList)
	}
	for i := range doc.Footers {
		stylizedBody(&doc.Footers[i], stylesList)
	}

	return &doc, nil
}

func parseContent(decoder *xml.Decoder, r *zip.ReadCloser, relationships Relationships) (Body, error) {
	var body Body
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return body, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "p" { // Paragraph
				var para Paragraph
				if err := decoder.DecodeElement(&para, &t); err != nil {
					return body, err
				}
				// Image handling is disabled as file output is not desired.
				body.Contents = append(body.Contents, ContentItem{Type: "paragraph", Value: para})
			} else if t.Name.Local == "tbl" { // Table
				var tbl Table
				if err := decoder.DecodeElement(&tbl, &t); err != nil {
					return body, err
				}
				body.Contents = append(body.Contents, ContentItem{Type: "table", Value: tbl})
			}
		}
	}
	return body, nil
}
