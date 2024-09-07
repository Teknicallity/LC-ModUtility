package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/inancgumus/screen"
	"lethalModUtility/internal/modList"
	"lethalModUtility/internal/pathUtil"
	"lethalModUtility/internal/zipUtil"
	"os"
	"path/filepath"
	"strings"
)

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
	fmt.Println("\t5. Install fresh from plugins.md file")
	fmt.Println("\tq. Quit and write to plugins.md")
	fmt.Println()
}

func main() {
	clearScreen()
	pluginsFile, err := findPluginsMd()
	if err != nil {
		fmt.Printf("Cannot find plugins.md file. Must be in local directory or ./Bepinex/")
		return
	}
	mods, err := modList.NewModListFromPluginsMd(pluginsFile)
	if err != nil {
		fmt.Printf("Mod List Error: %d", err)
		return
	}

	menu(mods)
}

func menu(m *modList.ModList) {
	defer func(m *modList.ModList) {
		err := m.WriteModsList()
		if err != nil {
			fmt.Printf("error writing mod list: %s", err)
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
			err := zipUtil.UnzipPack()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			successPrint()

		case "3":
			clearScreen()
			fmt.Println("Selected 3: Create new compressed modpack")
			err := zipUtil.ZipBepInEx(false)
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

			err = m.AddModFromUrl(link)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println()
			successPrint()

		case "5":
			clearScreen()
			fmt.Println("Selected 5: Installing fresh from plugins.md file")
			var withConfigChoice string
			fmt.Print("Would you like to keep the old config? (y/n): ")
			_, err = fmt.Scanln(&withConfigChoice)
			if err != nil {
				fmt.Println("Could not read selection")
				continue
			}

			var isInstallWithConfig bool
			if strings.ToLower(withConfigChoice)[0] == 'y' {
				isInstallWithConfig = true
			} else if strings.ToLower(withConfigChoice)[0] == 'n' {
				isInstallWithConfig = false
			} else {
				fmt.Println("Please enter 'y' or 'n'")
				continue
			}
			if isInstallWithConfig {
				err = pathUtil.MoveDir(filepath.Join("BepInEx", "config"), ".")
				if err != nil {
					fmt.Println("Could not move Bepinex config for safekeeping:", err)
					return
				}
			}

			// Remove previous backup folder
			if _, err = os.Stat("BepInExBackup.zip"); err == nil {
				err = os.RemoveAll("BepInExBackup.zip")
				if err != nil {
					fmt.Println("Could not remove backup:", err)
					return
				}
			}
			// Create backup
			err = zipUtil.ZipBepInEx(true)
			if err != nil {
				fmt.Println("Could not back up BepInEx folder:", err)
				return
			}
			err = os.RemoveAll("BepInEx")
			if err != nil {
				fmt.Println("Could not remove old BepInEx folder:", err)
				return
			}

			if isInstallWithConfig {
				err = pathUtil.MoveDir("config", "BepInEx")
				if err != nil {
					fmt.Println("Could not move Bepinex config into place:", err)
					return
				}
			}

			err = m.CleanInstallAllMods()
			if err != nil {
				fmt.Println("Could not install all mods:", err)
				return
			}

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

func findPluginsMd() (string, error) {
	paths := []string{
		filepath.Join("BepInEx", "plugins.md"),
		"plugins.md",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil // Return the valid path
		}
	}
	return "", os.ErrNotExist
}
