package docx_parser

import (
	"bytes"
	"github.com/mattn/go-runewidth"
	"strings"
)

func Table2markdown(table Table) string {
	// Align table width
	var tableWidth = 0
	for _, row := range table.Rows {
		if tableWidth < len(row.Cells) {
			tableWidth = len(row.Cells)
		}
	}
	// Align cell width
	alignArr := make([]int, tableWidth)
	for _, row := range table.Rows {
		for j, cell := range row.Cells {
			if alignArr[j] < vLen(strings.Join(cell.Texts, "")) {
				alignArr[j] = vLen(strings.Join(cell.Texts, ""))
			}
		}
	}

	var buffer bytes.Buffer
	// The first row is the header by default
	for i, row := range table.Rows {
		if i == 0 {
			widArr := make([]string, tableWidth) // For appending |--| table markers
			// Iterate through cells
			for j, cell := range row.Cells {
				cellStr := strings.Join(cell.Texts, "")
				var text = "| " + cellStr // Add
				buffer.WriteString(text)
				widDif := alignArr[j] - vLen(cellStr)
				if widDif > 0 {
					buffer.WriteString(strings.Repeat(" ", widDif)) // Add spaces
				} else {
					buffer.WriteString(" ")
				}

				widArr[j] = strings.Repeat("-", alignArr[j])
			}
			// Pad width
			for j := len(row.Cells); j < tableWidth; j++ {
				buffer.WriteString("| ")
				buffer.WriteString(strings.Repeat(" ", alignArr[j]))
				widArr[j] = strings.Repeat("-", alignArr[j])
			}
			buffer.WriteString("|\n")
			// -----------------------------
			// Appended after the first row.
			for j := range widArr {
				var text = "| " + widArr[j] + " "
				buffer.WriteString(text)
			}
			// ------------------------------
			// Always add this-----------------------------
			buffer.WriteString("|\n")
		} else { // Not the first row
			// Iterate through cells
			for j, cell := range row.Cells {
				cellStr := strings.Join(cell.Texts, "")
				var text = "| " + cellStr // Add
				buffer.WriteString(text)
				widDif := alignArr[j] - vLen(cellStr)
				if widDif > 0 {
					buffer.WriteString(strings.Repeat(" ", widDif)) // Add spaces
				} else {
					buffer.WriteString(" ")
				}
			}
			// Pad width
			for j := len(row.Cells); j < tableWidth; j++ {
				buffer.WriteString("| ")
				buffer.WriteString(strings.Repeat(" ", alignArr[j]))
			}
			buffer.WriteString("|\n")
		}
	}
	result := buffer.String()
	return result
}

func vLen(s string) int {
	// visual length
	return runewidth.StringWidth(s)
}
