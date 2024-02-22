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
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "modified_*.txt")
	if err != nil {
		return fmt.Errorf("error creating temporary file: %w", err)
	}
	defer tempFile.Close()

	// Create a scanner to read the input file line by line
	scanner := bufio.NewScanner(file)

	// Create a writer to write to the temporary file
	writer := bufio.NewWriter(tempFile)

	// Iterate over each line in the input file
	for scanner.Scan() {
		line := scanner.Text()

		// Modify the line if it meets the condition
		if strings.Contains(line, "-") {
			modAndVersion, modUrl := parseModLine(line)
			authorModVersion, downloadLink := scrapeWebsite(modUrl)

			localVersion := strings.Split(modAndVersion, "-")[1]
			remoteVersion := strings.Split(authorModVersion, "-")[1]
			if remoteVersion > localVersion {
				fmt.Printf("Updating: %s -> %s\n", modAndVersion, remoteVersion)
				err := downloadMod(downloadLink, authorModVersion)
				if err != nil {
					return err
				}

				line = fmt.Sprintf("- [%s](%s)", authorModVersion, modUrl) //proper
			}
		}
		// Write the modified line to the temporary file
		_, err := writer.WriteString(line + "\n")
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

func parseModLine(line string) (string, string) {
	var modVersion, modUrl string
	versionPattern := `\[(.*?)\]\(`
	vRe := regexp.MustCompile(versionPattern)
	vMatches := vRe.FindStringSubmatch(line)
	if len(vMatches) >= 2 {
		modVersion = vMatches[1]
	} else {
		fmt.Printf("ERR: No version pattern match found for %s\n", line)
	}
	urlPattern := `\]\((.*?)\)`
	urlRe := regexp.MustCompile(urlPattern)
	urlMatches := urlRe.FindStringSubmatch(line)
	if len(urlMatches) >= 2 {
		modUrl = urlMatches[1]
		fmt.Println("\tChecking:", modUrl)
	} else {
		fmt.Printf("ERR: No url pattern match found for %s\n", line)
	}

	return modVersion, modUrl
}
