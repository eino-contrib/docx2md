package docx_parser

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// extractImageFromDocx extracts and saves an image to a specified directory.
func extractImageFromDocx(r *zip.ReadCloser, rid string, relationships Relationships, outputDir string) (string, error) {
	// Image resource IDs usually start with "rId"
	var mediaFileName string
	for _, rel := range relationships.Relationship {
		if rid == rel.Id {
			mediaFileName = rel.Target
			break
		}
	}

	for _, f := range r.File {
		// Find the image file corresponding to the resource ID
		if strings.Contains(f.Name, mediaFileName) && strings.HasPrefix(f.Name, "word/media/") {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			defer func(rc io.ReadCloser) {
				err := rc.Close()
				if err != nil {
					// empty
				}
			}(rc)

			// Construct the output path for the image
			imagePath := filepath.Join(outputDir, filepath.Base(f.Name))

			// Save the image data to a file
			imageData, err := io.ReadAll(rc)
			if err != nil {
				return "", err
			}
			err = os.WriteFile(imagePath, imageData, os.ModePerm)
			if err != nil {
				return "", err
			}

			// Return the path of the image
			return imagePath, nil
		}
	}
	return "", fmt.Errorf("no image found for resource ID: %s", rid)
}
