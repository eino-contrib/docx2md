// -------------------------------------------------
// Package docx_parser
// Author: hanzhi
// Date: 2024/12/22
// -------------------------------------------------

package docx_parser

func stylizedBody(body *Body, styles *Styles) {
	// Create a map for easy lookup based on styles
	var styleFZMap map[string]int = make(map[string]int)
	for _, style := range styles.StyleList {
		if style.StyleId != "" {
			styleFZMap[style.StyleId] = style.FontSize.Value
		}
	}
	//fmt.Println(styleFZMap)
	// Iterate through the body, looking for paragraphs
	for i, content := range body.Contents {
		if content.Type == "paragraph" {
			paragraph := content.Value.(Paragraph)
			if paragraph.StyleId.Value != "" {
				if fontSize, exists := styleFZMap[paragraph.StyleId.Value]; exists {
					for j := range paragraph.Runs {
						paragraph.Runs[j].FontSize.Value = fontSize
					}
				}
			}
			// Write back to body.Contents[i]
			body.Contents[i].Value = paragraph
		}
	}

	//fmt.Println(styleFZMap)
}
