package scraper

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gocolly/colly"
	"lethalModUtility/internal/pathUtil"
	"lethalModUtility/internal/timeUtil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func DownloadMod(downloadUrl string, outputFileName string) (string, error) {
	pathToDownload := filepath.Join(pathUtil.GetDownloadFolderPath(), "LC_New_Mods", "zips")
	zipFilePath := filepath.Join(pathToDownload, outputFileName+".zip")

	c := colly.NewCollector(colly.MaxBodySize(100 * 1024 * 1024))
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Headers.Get("Content-Type"), "application/zip") {
			if _, err := os.Stat(pathToDownload); os.IsNotExist(err) {
				err = os.MkdirAll(pathToDownload, os.ModePerm)
				if err != nil {
					fmt.Println("cannot make new mod directory", err)
					return
				}
			}
			err := r.Save(zipFilePath)
			if err != nil {
				fmt.Println("Download zip file error:", err)
				return
			}
			clr := color.New(color.FgGreen)
			_, err = clr.Printf("Zip file downloaded successfully.")
			if err != nil {
				return
			}
		}
	})

	err := c.Visit(downloadUrl)
	if err != nil {
		return "", err
	}
	return zipFilePath, nil
}

func GetRemoteInfoFromUrl(url string) (RemoteInfo, error) {
	var downloadLink, modNameWithVersion, lastUpdatedHuman string
	// Create a new collector
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong when scraping: ", err)
	})
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	// Get most recent version
	c.OnHTML("meta[content]", func(e *colly.HTMLElement) {
		metaContent := e.Attr("content")
		if strings.Contains(metaContent, "https://gcdn.thunderstore.io/live/repository/icons/") {
			pattern := `icons\/((.*?)-((.*?)-([\d.]+)))\.png`
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(metaContent)
			// matches
			// 0 full icons/...
			// 1 author-mod-verstion
			// 2 author
			// 3 mod-version
			// 4 5 mod version

			if len(matches) >= 2 {
				modNameWithVersion = matches[3]
			} else {
				fmt.Printf("ERR: No match found for %s\n", metaContent)
			}
		}
	})

	// Get full zip download link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "https://thunderstore.io/package/download/") {
			downloadLink = link
		}
	})

	// Get last update time string
	c.OnHTML("table.table tr", func(e *colly.HTMLElement) {
		// Find the <td> element with the text "Last updated"
		if e.ChildText("td:first-child") == "Last updated" {
			// Extract and print the value in the adjacent <td>
			lastUpdatedHuman = e.ChildText("td:nth-child(2)")
		}
	})

	// On every <a> element inside <h5> tags, extract the href attribute and the text
	var dependencies []Dependency
	c.OnHTML("div.list-group-item div.media-body h5 a", func(e *colly.HTMLElement) {
		var modName, dependencyUrl string
		dependencyUrl = e.Attr("href") // "/c/lethal-company/p/BepInEx/BepInExPack/"
		modAuthorWithName := e.Text    // "BepInEx-BepInExPack"
		if dependencyUrl == "/c/lethal-company/p/BepInEx/BepInExPack/" {
			return
		}

		if !strings.Contains(dependencyUrl, "https://thunderstore.io/c/") {
			dependencyUrl = "https://thunderstore.io" + dependencyUrl
		}

		splitName := strings.Split(modAuthorWithName, "-")
		if len(splitName) > 0 {
			modName = strings.TrimSpace(splitName[len(splitName)-1])
		}
		dependencies = append(dependencies, Dependency{
			Name: modName,
			Url:  dependencyUrl,
		})
	})

	// Visit the URL
	err := c.Visit(url)
	if err != nil {
		return RemoteInfo{}, err
	}

	modVersion, err := parseVersionFromNameWithVersion(modNameWithVersion)
	if err != nil {
		return RemoteInfo{}, err
	}

	return RemoteInfo{
		ModVersion:               modVersion,
		DownloadUrl:              downloadLink,
		LastUpdatedHumanReadable: lastUpdatedHuman,
		LastUpdatedTime:          timeUtil.ParseTimeString(lastUpdatedHuman),
		ModNameWithVersion:       modNameWithVersion,
		Dependencies:             dependencies,
	}, nil
}

func parseVersionFromNameWithVersion(nameWithVersion string) (string, error) {
	versionPattern := `-([\d.]+)`
	versionRegx := regexp.MustCompile(versionPattern)
	remoteVersionMatches := versionRegx.FindStringSubmatch(nameWithVersion)
	if len(remoteVersionMatches) >= 2 {
		return remoteVersionMatches[1], nil
	}
	return "", fmt.Errorf("could not parse mod name with version: %s", nameWithVersion)
}
