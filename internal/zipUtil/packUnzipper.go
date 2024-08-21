package zipUtil

import (
	"archive/zip"
	"fmt"
	"lethalModUtility/internal/pathUtil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var directoryPath = "./"

func checkAndDeleteBepinexRelated() {
	bepinexPath := filepath.Join(directoryPath, "Bepinex")

	if _, err := os.Stat(bepinexPath); !os.IsNotExist(err) {
		os.RemoveAll(bepinexPath)
		fmt.Println("Bepinex folder deleted")
	} else {
		fmt.Println("Bepinex directory does not exist. Creating directory")
	}

	filesToRemove := []string{"winhttp.dll", "doorstop_config.ini"}
	for _, file := range filesToRemove {
		filePath := filepath.Join(directoryPath, file)
		if _, err := os.Stat(filePath); err == nil {
			os.Remove(filePath)
		}
	}
}

func unzipBepinex(downloadsPath, packName string) error {
	fmt.Println("Unzipping:", packName)
	zipInput := filepath.Join(downloadsPath, packName)
	zipReader, err := zip.OpenReader(zipInput)
	if err != nil {
		fmt.Println("Error opening zip file:", err)
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if !strings.HasPrefix(file.Name, "__MACOSX") {
			destPath := filepath.Join(directoryPath, file.Name)
			if file.FileInfo().IsDir() {
				err = os.MkdirAll(destPath, os.ModePerm)
				if err != nil {
					fmt.Println("cannot make new directory", err)
					return err
				}
			} else {
				err = extractFile(file, destPath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getBepinexPack(downloadsFolder string) string {
	pattern := regexp.MustCompile(`BepinExPack_v(\d+)`)
	files, err := os.ReadDir(downloadsFolder)
	if err != nil {
		fmt.Println("Error reading downloads folder:", err)
		return ""
	}

	var maxVersion int
	var maxVersionFileName string

	for _, file := range files {
		if pattern.MatchString(file.Name()) {
			fmt.Printf("\tFound: %s\n", file.Name())
			matches := pattern.FindStringSubmatch(file.Name())
			version := matches[1]

			if ver, err := strconv.Atoi(version); err == nil {
				if ver > maxVersion {
					maxVersion = ver
					maxVersionFileName = file.Name()
				}
			}
		}
	}
	return maxVersionFileName
}

func UnzipPack() error {

	checkAndDeleteBepinexRelated()
	downloadsPath := pathUtil.GetDownloadFolderPath()
	bepinexPackZip := getBepinexPack(downloadsPath)
	fmt.Println(downloadsPath)
	fmt.Println(bepinexPackZip)

	if bepinexPackZip == "" {
		fmt.Println("ERROR: Cannot find BepinexPack zip file")
		return os.ErrNotExist
	}

	err := unzipBepinex(downloadsPath, bepinexPackZip)
	if err != nil {
		return err
	}

	return nil
}
