package modList

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ModList struct {
	// file
	mods             []modEntry
	markDownFilePath string
}

func NewModList(modsListFilePath string) (*ModList, error) {
	mList := &ModList{
		mods:             nil,
		markDownFilePath: modsListFilePath,
	}

	err := mList.readModsMarkdownFile(modsListFilePath)
	if err != nil {
		return mList, fmt.Errorf("cannot read mods file: %w", err)
	}

	return mList, nil
}

func (m *ModList) readModsMarkdownFile(modsListFilePath string) error {
	file, err := os.Open(modsListFilePath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	modsSlice := make([]modEntry, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "-") {
			modEntry, err := newModEntryFromFileLine(line)
			if err != nil {
				return err
			}
			modsSlice = append(modsSlice, modEntry)
		}
	}
	m.mods = modsSlice

	return nil
}

func (m *ModList) AddMod(modUrl string) error {
	mod, err := newModEntryFromUrl(modUrl)
	if err != nil {
		return err
	}

	m.mods = append(m.mods, mod)

	return nil
}

func (m *ModList) UpdateAllMods() error {
	listLength := len(m.mods)
	fmt.Printf("%-9s %-25s %-18s %-18s %s\n", "Queue", "Mod Name", "Status", "Last Updated", "Action")

	for i, mod := range m.mods {
		sequence := fmt.Sprintf("[%d/%d]", i+1, listLength)
		fmt.Printf("%-9s ", sequence)

		err := mod.updateMod()
		if err != nil {
			return err
		}

		fmt.Println()
	}

	return nil
}

func (m *ModList) GrabModVersion() {

}

func (m *ModList) WriteModsList(outputDirector ...string) error {
	var err error
	var tempFile *os.File
	if len(outputDirector) != 0 {
		tempFile, err = os.CreateTemp(outputDirector[0], "modified_*.txt")
	} else {
		tempFile, err = os.CreateTemp("", "modified_*.txt")
	}
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	writer := bufio.NewWriter(tempFile)
	_, err = writer.WriteString("\n## Mod List" + "\n\n")
	if err != nil {
		return fmt.Errorf("error writing header to temporary file: %w", err)
	}

	for _, mod := range m.mods {
		line := mod.getMarkdownEntry()
		_, err = writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing mod to temporary file: %s:, %w", line, err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("error closing temporary file: %w", err)
	}

	if err := os.Rename(tempFile.Name(), m.markDownFilePath); err != nil {
		return fmt.Errorf("error renaming temporary file: %w", err)
	}

	return nil
}
