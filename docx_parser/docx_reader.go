package docx_parser

import (
	"archive/zip"
	"encoding/xml"
	"io"
	"strings"
)

// 修改后的读取和解析函数，将段落和表格混合存储
func ReadDocx(filePath string) (*Document, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var doc Document
	var relationships Relationships
	var stylesList *Styles

	// Pre-process relationships and styles
	for _, f := range r.File {
		if f.Name == "word/_rels/document.xml.rels" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			if err := xml.NewDecoder(rc).Decode(&relationships); err != nil {
				rc.Close()
				return nil, err
			}
			rc.Close()
		}
		if f.Name == "word/styles.xml" {
			stylesList, err = ReadStyle(r, filePath)
			if err != nil {
				return nil, err
			}
		}
	}

	// Process content files
	for _, f := range r.File {
		var bodyPart *Body
		isHeader := false
		isFooter := false

		if f.Name == "word/document.xml" {
			bodyPart = &doc.Body
		} else if strings.HasPrefix(f.Name, "word/header") && strings.HasSuffix(f.Name, ".xml") {
			isHeader = true
			doc.Headers = append(doc.Headers, Body{})
			bodyPart = &doc.Headers[len(doc.Headers)-1]
		} else if strings.HasPrefix(f.Name, "word/footer") && strings.HasSuffix(f.Name, ".xml") {
			isFooter = true
			doc.Footers = append(doc.Footers, Body{})
			bodyPart = &doc.Footers[len(doc.Footers)-1]
		} else {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		body, err := parseContent(xml.NewDecoder(rc), r, relationships)
		if err != nil {
			rc.Close()
			return nil, err
		}
		rc.Close()
		*bodyPart = body

		if isHeader || isFooter {
			// Optionally handle header/footer specific logic
		}
	}

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
