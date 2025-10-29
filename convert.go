package docx2md

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/eino-contrib/docx2md/docx_parser"
)

// Config controls the conversion process.
type Config struct {
	IncludeHeaders bool // whether to include headers in the parsed content
	IncludeFooters bool // whether to include footers in the parsed content
	IncludeTables  bool // whether to include table content
}

func DocxConvert(docPath string, config *Config) (map[string]string, error) { // Returns markdown string and error
	// --------------
	doc, err := docx_parser.ReadDocx(docPath)
	if err != nil {
		return nil, err
	}

	sections := make(map[string]string)
	var headerBuilder, bodyBuilder, tableBuilder, footerBuilder bytes.Buffer

	if config == nil {
		config = &Config{
			IncludeHeaders: true,
			IncludeFooters: true,
			IncludeTables:  true,
		}
	}

	// Process Headers
	if config.IncludeHeaders {
		for _, header := range doc.Headers {
			for _, content := range header.Contents {
				processContentItem(content, &headerBuilder, &tableBuilder, config)
			}
		}
	}

	// Process Body and Tables
	for _, content := range doc.Body.Contents {
		processContentItem(content, &bodyBuilder, &tableBuilder, config)
	}

	// Process Footers
	if config.IncludeFooters {
		for _, footer := range doc.Footers {
			for _, content := range footer.Contents {
				processContentItem(content, &footerBuilder, &tableBuilder, config)
			}
		}
	}

	if config.IncludeHeaders {
		sections["headers"] = headerBuilder.String()
	}
	sections["main"] = bodyBuilder.String()
	if config.IncludeTables {
		sections["tables"] = tableBuilder.String()
	}
	if config.IncludeFooters {
		sections["footers"] = footerBuilder.String()
	}

	return sections, nil
}

func processContentItem(content docx_parser.ContentItem, buffer, tableBuilder *bytes.Buffer, config *Config) {
	if content.Type == "paragraph" {
		var bufferPar bytes.Buffer
		para := content.Value.(docx_parser.Paragraph)
		numPr := para.NumPr

		var fontSizePar int
		// Handle paragraphs with no text runs.
		if len(para.Runs) == 0 {
			fontSizePar = 0 // Default font size for empty paragraphs.
		} else {
			// Initialize with the font size of the first run.
			fontSizePar = para.Runs[0].FontSize.Value
		}

		// Find the minimum font size and build the paragraph string.
		for _, run := range para.Runs {
			if run.FontSize.Value < fontSizePar {
				fontSizePar = run.FontSize.Value
			}
			for _, text := range run.Text {
				bufferPar.WriteString(text.Value)
			}
		}
		paragraphStr := bufferPar.String()
		paragraphStr = word2Heading(paragraphStr, fontSizePar, numPr)
		buffer.WriteString(paragraphStr)
	} else if config.IncludeTables && content.Type == "table" {
		table := content.Value.(docx_parser.Table)
		tableStr := docx_parser.Table2markdown(table)
		tableBuilder.WriteString(tableStr)
		tableBuilder.WriteString("\n")
	} else if content.Type == "image" {
		// Image processing is disabled
	}
	buffer.WriteString("\n")
}

func word2Heading(value string, fontSize int, numPr *bool) string {
	value = getTrimedStr(value)
	if len(value) == 0 {
		return "\n"
	}
	if numPr != nil || docx_parser.CheckString(value) {
		var maxHeadingLength = 45
		var h1 = 48
		var h2 = 36
		var h3 = 28
		var h4 = 24

		if h1 <= fontSize {
			return "# " + value
		} else if h2 <= fontSize {
			return "## " + value
		} else if h3 <= fontSize {
			return "### " + value
		} else if h4 <= fontSize && len(value) < maxHeadingLength {
			return "#### " + value
		} else {
			return value
		}

	} else {
		var maxHeadingLength = 15
		var h1 = 48
		var h2 = 36
		var h3 = 28

		if h1 <= fontSize {
			return "# " + value
		} else if h2 <= fontSize {
			return "## " + value
		} else if h3 <= fontSize && len(value) < maxHeadingLength {
			return "### " + value
		} else {
			return value
		}

	}

}

func getTrimedStr(s string) string {
	// Use TrimFunc to remove all leading/trailing whitespace from the string
	trimmed := strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	return trimmed
}
