// -------------------------------------------------
// Package docx_parser
// Author: hanzhi
// Date: 2024/12/22
// -------------------------------------------------

package docx_parser

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
)

func ReadStyle(r *zip.ReadCloser, filePath string) (*Styles, error) {
	// Find the styles.xml file
	var styleFileRels *zip.File
	for _, f := range r.File {
		if f.Name == "word/styles.xml" {
			styleFileRels = f
			break
		}
	}

	if styleFileRels == nil {
		return nil, fmt.Errorf("styles.xml not found in %s", filePath)
	}
	// Read the content of style.xml
	rcDFR, err := styleFileRels.Open()
	if err != nil {
		return nil, err
	}
	defer func(rc io.ReadCloser) {
		err := rc.Close()
		if err != nil {
			// empty
		}
	}(rcDFR)
	// Parse the XML
	var stylesList Styles
	err = xml.NewDecoder(rcDFR).Decode(&stylesList)
	if err != nil {
		return nil, err
	}

	return &stylesList, nil
}
