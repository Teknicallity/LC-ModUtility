package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"lethalModUtility/internal/modList"
	"lethalModUtility/internal/zipUtil"
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

func printInitialSelection() {
	fmt.Println("Options")
	fmt.Println("=============================")
	fmt.Println("\t1. Update mods")
	fmt.Println("\t2. Unzip pack from downloads")
	fmt.Println("\t3. Creating new compressed modpack")
	fmt.Println("\t4. Download new mod")
	fmt.Println("\tq. Quit")
	fmt.Println()
}

func main() {
	clearScreen()
	mods, err := modList.NewModList(".\\BepInEx\\plugins.md")
	if err != nil {
		fmt.Printf("Mod List Error: %d", err)
	}

	menu(mods)
}

func menu(m *modList.ModList) {
	defer func(m *modList.ModList) {
		err := m.WriteModsList()
		if err != nil {

		}
	}(m)

SelectionLoop:
	for {
		printInitialSelection()
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
			err := m.UpdateAllMods()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			successPrint("Checked All Mods")

		case "2":
			clearScreen()
			fmt.Println("Selected 2: Unziping pack from downloads")
			err := zipUtil.BuildPack()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			successPrint()

		case "3":
			clearScreen()
			fmt.Println("Selected 3: Create new compressed modpack")
			err := zipUtil.ZipBepinEx()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			successPrint()

		case "4":
			clearScreen()
			fmt.Println("Selected 4: Downloading new mod")
			link, err := getModLink()
			if err != nil {
				fmt.Printf("could not intake new mod link %d\n", err)
				return
			}

			err = m.AddMod(link)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			successPrint()

		case "q":
			//os.Exit(0)
			break SelectionLoop
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

func successPrint(message ...string) {
	if len(message) == 0 {
		message = append(message, "Success")
	}
	color.Green(message[0])
}
