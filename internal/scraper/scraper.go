package scraper

import (
	"fmt"
	"github.com/gocolly/colly"
	"lethalModUtility/internal/pathUtil"
	"os"
	"regexp"
	"strings"
)

func DownloadMod(downloadUrl string, outputFileName string) error {
	c := colly.NewCollector(colly.MaxBodySize(100 * 1024 * 1024))
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Headers.Get("Content-Type"), "application/zip") {
			pathToDownload := pathUtil.GetDownloadFolderPath() + "\\LC_New_Mods\\"
			if _, err := os.Stat(pathToDownload); os.IsNotExist(err) {
				err = os.Mkdir(pathToDownload, 0644)
				if err != nil {
					fmt.Println("cannot make new mod directory", err)
					return
				}
			}

			err := r.Save(pathToDownload + outputFileName + ".zip")
			if err != nil {
				fmt.Println("Download zip file error:", err)
				return
			}
			fmt.Printf("Zip file downloaded successfully.")
		}
	})

	err := c.Visit(downloadUrl)
	if err != nil {
		return err
	}
	return nil
}

func GetRemoteInfoFromUrl(url string) (RemoteInfo, error) {
	var downloadLink, modAndVersion, lastUpdated string
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
				modAndVersion = matches[3]
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
			lastUpdated = e.ChildText("td:nth-child(2)")
		}
	})

	// Visit the URL
	err := c.Visit(url)
	if err != nil {
		return RemoteInfo{}, err
	}

	modVersion, err := parseVersionFromNameWithVersion(modAndVersion)
	if err != nil {
		return RemoteInfo{}, err
	}

	return RemoteInfo{
		ModVersion:  modVersion,
		DownloadUrl: downloadLink,
		LastUpdated: lastUpdated,
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
