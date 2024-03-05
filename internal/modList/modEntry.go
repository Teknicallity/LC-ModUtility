package modList

import (
	"fmt"
	"lethalModUtility/internal/scraper"
	"regexp"
)

type modEntry struct {
	modName       string
	localVersion  string
	modUrl        string
	remoteVersion string
	downloadUrl   string
}

func newModEntryFromFileLine(line string) (modEntry, error) {
	mod := modEntry{}
	err := mod.parsePluginsLine(line)
	if err != nil {
		return modEntry{}, err
	}
	return mod, nil
}

func newModEntryFromUrl(modUrl string) (modEntry, error) {
	mod := modEntry{}
	remoteModAndVersion, downloadLink := scraper.GetIdAndDownloadLink(modUrl)
	err := scraper.DownloadMod(downloadLink, remoteModAndVersion)
	if err != nil {
		return mod, err
	}
	mod.fillInfoFromModAndVersion(remoteModAndVersion)
	mod.modUrl = modUrl

	return mod, nil
}

func (m *modEntry) parsePluginsLine(line string) error {
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

func (m *modEntry) fillRemoteVersionAndDownloadUrl() {
	var remoteVersion string

	remoteModAndVersion, downloadLink := scraper.GetIdAndDownloadLink(m.modUrl)
	versionPattern := `-([\d.]+)`
	versionRegx := regexp.MustCompile(versionPattern)
	remoteVersionMatches := versionRegx.FindStringSubmatch(remoteModAndVersion)
	if len(remoteVersionMatches) >= 2 {
		remoteVersion = remoteVersionMatches[1]
	}

	m.remoteVersion = remoteVersion
	m.downloadUrl = downloadLink
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
	m.localVersion = m.remoteVersion
	modAndVersion := fmt.Sprintf(m.modName + "-" + m.localVersion)
	err := scraper.DownloadMod(m.downloadUrl, modAndVersion)
	if err != nil {
		return "", err
	}

	return modAndVersion, nil
}
