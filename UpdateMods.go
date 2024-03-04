package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func UpdateMods() error {
	err := modifyFile(".\\BepInEx\\plugins.md")
	if err != nil {
		return err
	}
	return nil
}

func modifyFile(inputFilePath string) error {
	// Open the input file
	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "modified_*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {

		}
	}(tempFile)

	// Create a scanner to read the input file line by line
	scanner := bufio.NewScanner(file)

	// Create a writer to write to the temporary file
	writer := bufio.NewWriter(tempFile)

	// Iterate over each line in the modlist file
	for scanner.Scan() {
		oldLine := scanner.Text()
		var newLine string

		// Modify the line if it is part of the modlist
		if strings.Contains(oldLine, "-") {
			newLine, err = processPluginsLine(oldLine)
			if err != nil {
				return err
			}
		} else { // if it is just the header or blank line, pass straight to new file
			newLine = oldLine
		}
		// Write the modified line to the temporary file
		_, err := writer.WriteString(newLine + "\n")
		if err != nil {
			return fmt.Errorf("error writing to temporary file: %w", err)
		}
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	// Flush the writer to ensure all buffered data is written to the file
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	// Close both the original file and the temporary file
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing original file: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("error closing temporary file: %w", err)
	}

	// Rename the temporary file to the original file name
	if err := os.Rename(tempFile.Name(), inputFilePath); err != nil {
		return fmt.Errorf("error renaming temporary file: %w", err)
	}

	return nil
}

func processPluginsLine(line string) (string, error) {
	var newLine string
	modAndVersion, modUrl := parsePluginsLine(line)
	authorModVersion, downloadLink := scrapeWebsite(modUrl)

	localVersionNumber := strings.Split(modAndVersion, "-")[1]
	remoteVersionNumber := strings.Split(authorModVersion, "-")[1]
	if remoteVersionNumber > localVersionNumber {
		fmt.Printf("Updating: %s -> %s\n", modAndVersion, remoteVersionNumber)
		err := downloadMod(downloadLink, authorModVersion)
		if err != nil {
			return "", err
		}

		newLine = fmt.Sprintf("- [%s](%s)", authorModVersion, modUrl) //proper
	} else {
		newLine = line
	}
	return newLine, nil
}

func parsePluginsLine(line string) (string, string) {
	var modVersion, modUrl string

	versionPattern := `\[(.*?)\]\(`
	versionRegx := regexp.MustCompile(versionPattern)
	versionMatches := versionRegx.FindStringSubmatch(line)
	if len(versionMatches) >= 2 {
		modVersion = versionMatches[1]
	} else {
		fmt.Printf("ERR: No version pattern match found for %s\n", line)
	}

	urlPattern := `\]\((.*?)\)`
	urlRegx := regexp.MustCompile(urlPattern)
	urlMatches := urlRegx.FindStringSubmatch(line)
	if len(urlMatches) >= 2 {
		modUrl = urlMatches[1]
		fmt.Println("\tChecking:", modUrl)
	} else {
		fmt.Printf("ERR: No url pattern match found for %s\n", line)
	}

	return modVersion, modUrl
}
