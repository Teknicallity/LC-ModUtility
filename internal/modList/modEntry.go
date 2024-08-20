package modList

import (
	"fmt"
	"github.com/fatih/color"
	"lethalModUtility/internal/scraper"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

type modEntry struct {
	modName      string
	localVersion string
	modUrl       string
	//remoteVersion string
	//downloadUrl   string
	//lastUpdated   string
	remoteInfo scraper.RemoteInfo
}

func newModEntryFromPluginsMdLine(line string) (modEntry, error) {
	mod := modEntry{}
	err := mod.parsePluginsMdLine(line)
	if err != nil {
		return modEntry{}, err
	}
	return mod, nil
}

func newModEntryFromUrl(modUrl string) (modEntry, error) {
	mod := modEntry{}
	remoteInfo, err := scraper.GetRemoteInfoFromUrl(modUrl)
	if err != nil {
		return mod, err
	}
	mod.fillInfoFromModAndVersion(remoteInfo.ModNameWithVersion)
	mod.modUrl = modUrl
	mod.remoteInfo = remoteInfo

	return mod, nil
}

func (m *modEntry) parsePluginsMdLine(line string) error {
	var mod, version, modUrl string

	//versionPattern := `\[(.*?)\]\(`
	versionPattern := `\[(.*?)-((\d{1,4}\.){2,3}\d{1,4})\]\(` // parse [mod-0.0.0](
	versionRegx := regexp.MustCompile(versionPattern)
	versionMatches := versionRegx.FindStringSubmatch(line)
	if len(versionMatches) >= 2 {
		mod = versionMatches[1]
		version = versionMatches[2]
	} else {
		return fmt.Errorf("No version pattern match found: %s\n", line)
	}

	urlPattern := `\]\((.*?)\)` // parse ](url)
	urlRegx := regexp.MustCompile(urlPattern)
	urlMatches := urlRegx.FindStringSubmatch(line)
	if len(urlMatches) >= 2 {
		modUrl = urlMatches[1]
		fmt.Println("\tChecking:", modUrl)

	} else {
		return fmt.Errorf("No url pattern match found: %s\n", line)
	}

	m.modName = mod
	m.localVersion = version
	m.modUrl = modUrl

	return nil
}

func (m *modEntry) fillRemoteInfo() {

	//remoteModAndVersion, downloadLink := scraper.GetRemoteInfoFromUrl(m.modUrl)
	remoteInfo, err := scraper.GetRemoteInfoFromUrl(m.modUrl)
	if err != nil {
		return
	}

	m.remoteInfo = remoteInfo
}

func (m *modEntry) fillInfoFromModAndVersion(modAndVersion string) {
	pattern := `(.*?)-((\d{1,4}\.){2,3}\d{1,4})` // parse mod-0.0.0
	regx := regexp.MustCompile(pattern)
	matches := regx.FindStringSubmatch(modAndVersion)
	if len(matches) >= 2 {
		m.modName = matches[1]
		m.localVersion = matches[2]
	}
}

func (m *modEntry) getMarkdownEntry() string {
	return fmt.Sprintf("- [%s-%s](%s)\n", m.modName, m.localVersion, m.modUrl)
}

func (m *modEntry) downloadMod() (string, error) {
	m.localVersion = m.remoteInfo.ModVersion
	modNameWithVersion := fmt.Sprintf(m.modName + "-" + m.localVersion)
	zipFilePath, err := scraper.DownloadMod(m.remoteInfo.DownloadUrl, modNameWithVersion)
	if err != nil {
		return "", err
	}

	return zipFilePath, nil
}

func (m *modEntry) updateMod() (string, error) {
	modNameString := fmt.Sprintf("%s:", m.modName)
	fmt.Printf("%-25s ", modNameString)

	if &m.remoteInfo == nil {
		return "", fmt.Errorf("remote info is nil")
	}

	if m.remoteInfo.ModVersion > m.localVersion {
		versionUpgrade := fmt.Sprintf("%s -> %s", m.localVersion, m.remoteInfo.ModVersion)
		fmt.Printf("%-18s ", versionUpgrade)
		m.printLastUpdatedString()
		zipFilePath, err := m.downloadMod()
		if err != nil {
			return "", fmt.Errorf("could not download %s: %w\n", filepath.Base(zipFilePath), err)
		}
		return zipFilePath, nil
	}
	upToDate := fmt.Sprintf("Up to date.")
	fmt.Printf("%-18s ", upToDate)
	m.printLastUpdatedString()

	return "", nil
}

func (m *modEntry) printLastUpdatedString() {
	var c *color.Color
	now := time.Now()

	switch {
	case m.remoteInfo.LastUpdatedTime.After(now.AddDate(0, 0, -1)):
		c = color.New(color.FgHiCyan)
	case m.remoteInfo.LastUpdatedTime.After(now.AddDate(0, 0, -8)):
		c = color.New(color.FgCyan)
	case m.remoteInfo.LastUpdatedTime.After(now.AddDate(0, -1, -1)):
		c = color.New(color.FgHiBlue)
	case m.remoteInfo.LastUpdatedTime.After(now.AddDate(0, -3, -1)):
		c = color.New(color.FgHiGreen)
	case m.remoteInfo.LastUpdatedTime.After(now.AddDate(0, -6, -1)):
		c = color.New(color.FgHiYellow)
	case m.remoteInfo.LastUpdatedTime.After(now.AddDate(-1, 0, -1)):
		c = color.New(color.FgYellow)
	case m.remoteInfo.LastUpdatedTime.After(now.AddDate(-2, 0, -1)):
		c = color.New(color.FgHiRed)
	default:
		c = color.New(color.FgRed)
	}

	whitelist := []string{
		"bepinex",
	}
	if slices.Contains(whitelist, strings.ToLower(m.modName)) {
		c = color.New(color.FgHiBlack)
	}

	_, err := c.Printf("%-16s ", m.remoteInfo.LastUpdatedHumanReadable)
	if err != nil {
		return
	}
}
