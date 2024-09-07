package zipUtil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func isFileAllowed(fileName string) bool {
	blacklistedFilesNames := []string{
		"readme.md",
		"manifest.json",
		"license",
		"icon.png",
		"changelog.md",
	}
	for _, blacklistedFile := range blacklistedFilesNames {
		if blacklistedFile == strings.ToLower(fileName) {
			return false
		}
	}
	return true
}

func extractFile(file *zip.File, destPath string) error {
	if file.FileInfo().IsDir() {
		return nil
	}

	if !isFileAllowed(filepath.Base(destPath)) {
		return nil
	}

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

func UnzipFile(zipFilePath string, destDirectory string) (string, error) {
	// destDirectory: user/downloads/LC_New_Mods/
	// zipFilePath: user/downloads/LC_New_Mods/zip/foo.zip
	zipFileName := filepath.Base(zipFilePath)                                        // foo.zip
	unzippedFolderName := strings.TrimSuffix(zipFileName, filepath.Ext(zipFileName)) // foo

	unzippedFolderPath := filepath.Join(destDirectory, unzippedFolderName) // user/downloads/LC_New_Mods//foo

	zipReader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("Error opening zip file, %s: %w\n", unzippedFolderPath, err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "__MACOSX") {
			return "", fmt.Errorf("unsupported mac file type: %s", file.Name)
		}

		destinationFilePath := filepath.Join(unzippedFolderPath, file.Name) //print
		err = extractFile(file, destinationFilePath)
		if err != nil {
			return "", err
		}
	}
	return unzippedFolderPath, nil
}
