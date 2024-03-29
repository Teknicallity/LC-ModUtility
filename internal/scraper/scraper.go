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
			fmt.Println("Zip file downloaded successfully.")
		}
	})

	err := c.Visit(downloadUrl)
	if err != nil {
		return err
	}
	return nil
}

func GetIdAndDownloadLink(url string) (string, string) {
	var downloadLink, modAndVersion string
	// Create a new collector
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong when scraping: ", err)
	})
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	// Set up event handlers
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

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "https://thunderstore.io/package/download/") {
			downloadLink = link
		}
	})

	// Visit the URL
	err := c.Visit(url)
	if err != nil {
		return "", ""
	}

	return modAndVersion, downloadLink
}
