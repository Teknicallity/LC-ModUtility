package zipUtil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractFile(file *zip.File, destPath string) error {
	err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
	if err != nil {
		fmt.Println("cannot make new directory", err)
		return err
	}

	extractedFile, err := file.Open()
	if err != nil {
		fmt.Println("Error extracting file:", err)
		return err
	}
	defer extractedFile.Close()

	fileBytes, err := io.ReadAll(extractedFile)
	if err != nil {
		fmt.Println("Error reading file contents:", err)
		return err
	}

	err = os.WriteFile(destPath, fileBytes, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}
	return nil
}

func unzipFile(zipFilePath string, destDirectory string) error {
	// destDirectory: user/downloads/LC_New_Mods/
	// zipFilePath: user/downloads/LC_New_Mods/zip/foo.zip
	zipFileName := filepath.Base(zipFilePath) // foo.zip
	unzippedFolderName := strings.TrimSuffix(zipFileName, filepath.Ext(zipFileName))

	unzippedFolderPath := filepath.Join(destDirectory, unzippedFolderName)

	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "__MACOSX") {
			return fmt.Errorf("unsupported mac file type: %s", file.Name)
		}

		destinationFilePath := filepath.Join(unzippedFolderPath, file.Name)
		err = extractFile(file, destinationFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}
