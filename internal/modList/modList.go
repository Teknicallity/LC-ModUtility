package modList

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"lethalModUtility/internal/modInstaller"
	"lethalModUtility/internal/pathUtil"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

type ModList struct {
	// file
	mods             []modEntry
	markDownFilePath string
	updatedMods      []string
}

func NewModListFromPluginsMd(modsListFilePath string) (*ModList, error) {
	mList := &ModList{
		mods:             nil,
		markDownFilePath: modsListFilePath,
	}

	err := mList.readModsMarkdownFile(modsListFilePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read mods file: %w", err)
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
		if strings.HasPrefix(line, "-") {
			modEntry, err := newModEntryFromPluginsMdLine(line)
			if err != nil {
				return err
			}
			//modsSlice = append(modsSlice, modEntry)

			modsSlice = m.AddModEntryToList(modsSlice, modEntry)
		}
	}
	m.mods = modsSlice

	return nil
}

func (m *ModList) AddModEntryToList(modsSlice []modEntry, modEntry modEntry) []modEntry {
	index := sort.Search(len(modsSlice), func(i int) bool {
		return modsSlice[i].modName > modEntry.modName
	})

	// result = slices.Insert(slice, index, value)
	// Insert the modEntry at the found index
	modsSlice = slices.Insert(modsSlice, index, modEntry)
	return modsSlice
}

//func (m *ModList) AddModEntryToList(modEntry *modEntry) {
//
//}

func (m *ModList) AddModFromUrl(modUrl string) error {
	mod, err := newModEntryFromUrl(modUrl)
	if err != nil {
		return err
	}

	for i := range m.mods {
		if m.mods[i].modName == mod.modName {
			return nil
		}
	}

	zipFilePath, err := mod.downloadMod()
	if err != nil {
		return fmt.Errorf("could not download %s: %w\n", filepath.Base(zipFilePath), err)
	}

	if zipFilePath != "" {
		err = modInstaller.InstallModFromZip(zipFilePath)
		if err != nil {
			return fmt.Errorf("could not install mod from zip file: %w\n", err)
		}
	}

	m.mods = append(m.mods, mod)

	fmt.Println()
	err = m.handleDependencies([]modEntry{mod})
	if err != nil {
		return fmt.Errorf("could not handle dependencies: %w\n", err)
	}

	return nil
}

func (m *ModList) UpdateAllMods() error {
	var modsThatNeedDependencyInstalls []modEntry
	listLength := len(m.mods)
	fmt.Printf("%-9s %-25s %-18s %-18s %s\n", "Queue", "Mod Name", "Status", "Last Updated", "Action")

	for i := range m.mods {
		sequence := fmt.Sprintf("[%d/%d]", i+1, listLength)
		fmt.Printf("%-9s ", sequence)
		modNameString := fmt.Sprintf("%s:", m.mods[i].modName)
		fmt.Printf("%-25s ", modNameString)

		m.mods[i].fillRemoteInfo()

		unfulfilledDependencies := m.mods[i].getUnfulfilledDependencies(m)
		if unfulfilledDependencies != nil {
			modsThatNeedDependencyInstalls = append(modsThatNeedDependencyInstalls, m.mods[i])
		}

		zipFilePath, err := m.mods[i].updateMod()
		if err != nil {
			return err
		}

		if zipFilePath != "" {
			err = modInstaller.InstallModFromZip(zipFilePath)
			if err != nil {
				return err
			}
		}
		fmt.Println()
	}

	if len(modsThatNeedDependencyInstalls) == 0 {
		return nil
	}

	fmt.Println()
	fmt.Printf("%-9s %-25s %-25s %s\n", "", "Dependent", "Dependency Name", "Action")
	if err := m.handleDependencies(modsThatNeedDependencyInstalls); err != nil {
		return err
	}

	return nil
}

func (m *ModList) handleDependencies(modEntriesInNeed []modEntry) error {
	for _, mEntry := range modEntriesInNeed {

		for i, dependency := range mEntry.remoteInfo.Dependencies {
			sequence := fmt.Sprintf("[%d/%d]", i+1, len(mEntry.remoteInfo.Dependencies))
			fmt.Printf("%-9s%-25s ", sequence, mEntry.modName)
			dependencyName := fmt.Sprintf("%s:", dependency.Name)
			fmt.Printf("%-25s ", dependencyName)

			if err := m.AddModFromUrl(dependency.Url); err != nil {
				return fmt.Errorf("error adding dependency from %s: %w", dependency.Url, err)
			}
		}
	}
	fmt.Println()
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

	err = pathUtil.MoveFile(tempFile.Name(), m.markDownFilePath)
	if err != nil {
		return err
	}

	c := color.New(color.FgGreen)
	_, err = c.Printf("Wrote to %s\n", m.markDownFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (m *ModList) doesModEntryExist(modName string) *modEntry {
	index := sort.Search(len(m.mods), func(i int) bool {
		return m.mods[i].modName >= modName
	})

	if index < len(m.mods) && m.mods[index].modName == modName {
		return &m.mods[index]
	}
	return nil
}

// CleanInstallAllMods Downloads and installs the latest version of all mods regardless plugins file
// Extremely Similar to UpdateAllMods. Refactor needed
func (m *ModList) CleanInstallAllMods() error {
	var modsThatNeedDependencyInstalls []modEntry
	listLength := len(m.mods)
	fmt.Printf("%-9s %-25s %-18s %-18s %s\n", "Queue", "Mod Name", "Status", "Last Updated", "Action")

	for i := range m.mods {
		sequence := fmt.Sprintf("[%d/%d]", i+1, listLength)
		fmt.Printf("%-9s ", sequence)
		modNameString := fmt.Sprintf("%s:", m.mods[i].modName)
		fmt.Printf("%-25s ", modNameString)

		m.mods[i].fillRemoteInfo()
		fmt.Printf("%-18s", m.mods[i].remoteInfo.ModVersion)

		unfulfilledDependencies := m.mods[i].getUnfulfilledDependencies(m)
		if unfulfilledDependencies != nil {
			modsThatNeedDependencyInstalls = append(modsThatNeedDependencyInstalls, m.mods[i])
		}

		m.mods[i].printLastUpdatedString()
		zipFilePath, err := m.mods[i].downloadMod()
		if err != nil {
			return fmt.Errorf("could not download %s: %w\n", filepath.Base(zipFilePath), err)
		}

		if zipFilePath != "" {
			err = modInstaller.InstallModFromZip(zipFilePath)
			if err != nil {
				return err
			}
		}
		fmt.Println()
	}

	if len(modsThatNeedDependencyInstalls) == 0 {
		return nil
	}

	fmt.Println()
	fmt.Printf("%-9s %-25s %-25s %s\n", "", "Dependent", "Dependency Name", "Action")
	if err := m.handleDependencies(modsThatNeedDependencyInstalls); err != nil {
		return err
	}

	return nil
}
