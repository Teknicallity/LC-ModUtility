package main

import (
	"fmt"
	"github.com/inancgumus/screen"
	"os"
	"path/filepath"
)

func getDownloadFolderPath() string {
	userProfile := os.Getenv("USERPROFILE")
	downloadsFolder := filepath.Join(userProfile, "Downloads")
	return downloadsFolder
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
	fmt.Println("\t2. Unzip pack from downloads")
	fmt.Println("\t2. Creating new compressed modpack")
	fmt.Println("\t4. Download new mod")
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
			err := UpdateMods()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println("Checked All Mods")

		case "2":
			clearScreen()
			fmt.Println("Selected 2: Unziping pack from downloads")
			err := BuildPack()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

		case "3":
			clearScreen()
			fmt.Println("Selected 3: Create new compressed modpack")
			err := ZipBepinEx()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

		case "4":
			clearScreen()
			fmt.Println("Selected 4: Downloading new mod")
			link, err := getModLink()
			if err != nil {
				fmt.Printf("could not intake new mod link %d\n", err)
			}
			err = ScrapeDownloadMod(link)
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
