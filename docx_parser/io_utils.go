package docx_parser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateMdDir(documentPath string, outputDir string, suffix string) (string, string, error) {
	// Get the file name and change the extension
	docxName := filepath.Base(documentPath)
	mdName := strings.TrimSuffix(docxName, suffix)
	mdName = mdName + ".md"

	// ---------------------------------
	// Check if outputDir exists, if not, return an error
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		fmt.Printf("path does not exist: %s\n", outputDir)
		return "", "", err
	}
	// --------------------------------
	// Create a path with uuid
	//uuidStr := uuid.New().String()
	//mdDirPath := filepath.Join(outputDir, uuidStr)
	mdDirPath := outputDir
	mdPath := filepath.Join(mdDirPath, mdName)
	// Check if the path exists, if not, create it
	err := os.MkdirAll(mdDirPath, 0755)
	if err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		return "", "", err
	}
	return mdPath, mdDirPath, err
}

func SaveFile(filePath string, mdStr string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			// Do nothing
		}
	}(file)
	// Write content
	_, err = file.Write([]byte(mdStr))
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
		return err
	}
	//
	return err
}
