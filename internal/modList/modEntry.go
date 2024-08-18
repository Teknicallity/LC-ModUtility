package modList

import (
	"fmt"
	"lethalModUtility/internal/scraper"
	"regexp"
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

func newModEntryFromFileLine(line string) (modEntry, error) {
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
	err = scraper.DownloadMod(remoteInfo.ModVersion, remoteInfo.ModVersion)
	if err != nil {
		return mod, err
	}
	mod.fillInfoFromModAndVersion(remoteInfo.ModVersion)
	mod.modUrl = modUrl

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
	modAndVersion := fmt.Sprintf(m.modName + "-" + m.localVersion)
	err := scraper.DownloadMod(m.remoteInfo.DownloadUrl, modAndVersion)
	if err != nil {
		return "", err
	}

	return modAndVersion, nil
}

func (m *modEntry) updateMod() error {
	modNameString := fmt.Sprintf("%s:", m.modName)
	fmt.Printf("%-25s ", modNameString)

	m.fillRemoteInfo()
	if m.remoteInfo.ModVersion > m.localVersion {
		versionUpgrade := fmt.Sprintf("%s -> %s", m.localVersion, m.remoteInfo.ModVersion)
		fmt.Printf("%-18s %s\t", versionUpgrade, m.remoteInfo.LastUpdated)
		modAndVersion, err := m.downloadMod()
		if err != nil {
			return fmt.Errorf("could not download %s: %w\n", modAndVersion, err)
		}
	} else {
		upToDate := fmt.Sprintf("Up to date.")
		fmt.Printf("%-18s %s", upToDate, m.remoteInfo.LastUpdated)
	}
	return nil
}
