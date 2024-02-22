package main

import (
	"bufio"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/inancgumus/screen"
	"os"
	"path/filepath"
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

func updateMods() error {
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
				fmt.Printf("Updating: %s\n", modAndVersion)
				err := downloadMod(downloadLink, modAndVersion)
				if err != nil {
					return err
				}

				line = fmt.Sprintf("- [%s](%s)", authorModVersion, modUrl)
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

func scrapeDownloadMod(modUrl string) error {
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

func getDownloadFolderPath() string {
	userProfile := os.Getenv("USERPROFILE")
	downloadsFolder := filepath.Join(userProfile, "Downloads")
	return downloadsFolder
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

func getModLink() (string, error) {
	var url string
	fmt.Printf("Enter new mod link: ")
	_, err := fmt.Scanln(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func initialSelection() {
	fmt.Println("Options")
	fmt.Println("=============================")
	fmt.Println("\t1. Update mods")
	fmt.Println("\t2. Build pack")
	fmt.Println("\t3. Download new mod")
	fmt.Println("\tq. Quit")
	fmt.Println()
}

func main() {
	clearScreen()
	for {
		initialSelection()
		var selection string
		fmt.Print("Please make selection by number: ")
		_, err := fmt.Scanln(&selection)
		if err != nil {
			fmt.Println("error read")
			return
		}
		switch selection {
		case "1":
			clearScreen()
			fmt.Println("Selected 1: Updating mods")
			err := updateMods()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Checked All Mods")

		case "2":
			clearScreen()
			fmt.Println("Selected 2: Building Pack")
			err := BuildPack()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

		case "3":
			clearScreen()
			fmt.Println("Selected 3: Downloading new mod")
			link, err := getModLink()
			if err != nil {
				fmt.Printf("could not intake new mod link %d\n", err)
			}
			err = scrapeDownloadMod(link)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

		case "q":
			os.Exit(0)
		default:
			clearScreen()
			fmt.Println("Could not understand selection")
		}
	}
}

func clearScreen() {
	screen.Clear()
	screen.MoveTopLeft()
}
