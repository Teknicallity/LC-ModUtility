package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"os"
	"regexp"
	"strings"
)

func scrapeWebsite(url string) (string, string) {
	var downloadLink, authorModVersion string

	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong when scraping: ", err)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	// parses auther-mod-version
	c.OnHTML("meta[content]", func(e *colly.HTMLElement) {
		metaContent := e.Attr("content")
		if strings.Contains(metaContent, "https://gcdn.thunderstore.io/live/repository/icons/") {
			//fmt.Printf("Link with Mod auther and version found: %s\n", metaContent)

			pattern := `icons\/(.*?)-(.*?)\.png`
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(metaContent)

			if len(matches) >= 2 {
				authorModVersion = matches[2]
			} else {
				fmt.Printf("ERR: No match found for %s\n", metaContent)
			}
		}
	})

	// gets the download link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "https://thunderstore.io/package/download/") {
			//fmt.Printf("Mod download found:%s\n", link)
			downloadLink = link
		}
	})

	c.Visit(url)

	return authorModVersion, downloadLink
}

func ScrapeDownloadMod(modUrl string) error {
	modAndVersion, downloadLink := scrapeWebsite(modUrl)
	newModLine := fmt.Sprintf("- [%s](%s)\n", modAndVersion, modUrl)

	err := downloadMod(downloadLink, modAndVersion)
	if err != nil {
		return err
	}

	err = appendToFile(newModLine)
	if err != nil {
		return err
	}

	return nil
}
func downloadMod(downloadUrl string, filename string) error {
	c := colly.NewCollector(colly.MaxBodySize(100 * 1024 * 1024))
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Headers.Get("Content-Type"), "application/zip") {
			pathToDownload := getDownloadFolderPath() + "\\LC_New_Mods\\"
			if _, err := os.Stat(pathToDownload); err != nil {
				if os.IsNotExist(err) {
					err = os.Mkdir(pathToDownload, 0644)
					if err != nil {
						fmt.Println("Make directory error", err)
						return
					}
				}
			}

			err := r.Save(pathToDownload + filename + ".zip")
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

func appendToFile(line string) error {

	file, err := os.OpenFile(".\\BepInEx\\plugins.md", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to the end of the file
	_, err = file.WriteString(line)
	if err != nil {
		return err
	}
	fmt.Println("Data appended successfully.")
	return nil
}
