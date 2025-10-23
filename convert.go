package docx2md

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/eino-contrib/docx2md/docx_parser"
)

func DocxConvert(docPath string) (string, error) { // Returns markdown string and error
	// --------------
	doc, err := docx_parser.ReadDocx(docPath)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	// Process Headers
	for _, header := range doc.Headers {
		for _, content := range header.Contents {
			processContentItem(content, &buffer)
		}
	}

	// Process Body
	for _, content := range doc.Body.Contents {
		processContentItem(content, &buffer)
	}

	// Process Footers
	for _, footer := range doc.Footers {
		for _, content := range footer.Contents {
			processContentItem(content, &buffer)
		}
	}

	return buffer.String(), nil
}

func processContentItem(content docx_parser.ContentItem, buffer *bytes.Buffer) {
	if content.Type == "paragraph" {
		var bufferPar bytes.Buffer
		var fontSizePar int = 48
		para := content.Value.(docx_parser.Paragraph)
		numPr := para.NumPr
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
	} else if content.Type == "table" {
		table := content.Value.(docx_parser.Table)
		tableStr := docx_parser.Table2markdown(table)
		buffer.WriteString(tableStr)
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
