package docx_parser

import "encoding/xml"

type ContentItem struct {
	Type  string      // "paragraph", "table", "image"
	Value interface{} // Paragraph, Table, or image path (string)
}

type Document struct {
	Body    Body   `xml:"body"`
	Headers []Body // for header*.xml
	Footers []Body // for footer*.xml
}

type Body struct {
	Contents []ContentItem
}

type Paragraph struct { // Paragraph type w:p
	Runs    []Run  `xml:"r"`          // A paragraph contains multiple runs
	NumPr   *bool  `xml:"pPr>numPr"`  // Check for numbering information
	StyleId PStyle `xml:"pPr>pStyle"` // Paragraph style
}

type PStyle struct { // Paragraph style
	Value string `xml:"val,attr"`
}

type Run struct { // Text run, may contain text or images
	FontSize FontSize `xml:"rPr>sz"`
	//FontBold    *bool    `xml:"rPr>b"`        // Removed due to continuous bold issues
	//FontIncline *bool    `xml:"rPr>i"`
	Text    []Text   `xml:"t"`
	Drawing *Drawing `xml:"drawing,omitempty"` // May contain an image
}

type Drawing struct { // Image nested in <w:drawing>
	Blip Blip `xml:"inline>graphic>graphicData>pic>blipFill>blip"`
}

type Blip struct {
	Embed string `xml:"embed,attr"`
}

type Table struct {
	XMLName xml.Name `xml:"tbl"`
	Rows    []Row    `xml:"tr"`
}

type Row struct {
	Cells []Cell `xml:"tc"`
}

type Cell struct {
	Texts []string `xml:"p>r>t"`
}

type Text struct {
	Value string `xml:",chardata"`
}

// FontSize -------------------------------
type FontSize struct {
	Value int `xml:"val,attr"`
}

// Relationships -------------------------------------
// Define struct to match XML format
type Relationships struct {
	XMLName      xml.Name       `xml:"http://schemas.openxmlformats.org/package/2006/relationships Relationships"`
	Relationship []Relationship `xml:"http://schemas.openxmlformats.org/package/2006/relationships Relationship"`
}

type Relationship struct {
	Id     string `xml:"Id,attr"`
	Type   string `xml:"Type,attr"`
	Target string `xml:"Target,attr"`
}

/* -------------------------------------------------------------- */

// Styles stylesheet
type Styles struct {
	XMLName   xml.Name
	StyleList []Style `xml:"style"`
}

type Style struct {
	Name     Name     `xml:"name"`
	StyleId  string   `xml:"styleId,attr"`
	FontSize FontSize `xml:"rPr>sz"`
}

type Name struct {
	Value string `xml:"val,attr"`
}
