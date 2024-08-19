package modInstaller

import (
	"fmt"
	"github.com/fatih/color"
	"lethalModUtility/internal/pathUtil"
	"lethalModUtility/internal/zipUtil"
	"os"
	"path/filepath"
)

func InstallMod(zipFilePath string) error {
	unzippedFolderPath, err := unzipMod(zipFilePath)
	if err != nil {
		return err
	}

	err = moveModFiles(unzippedFolderPath)
	if err != nil {
		return err
	}

	return nil
}

func moveModFiles(unzippedFolderPath string) error {
	fileInfo, err := os.ReadDir(unzippedFolderPath)
	if err != nil {
		fmt.Printf("error reading mod files: %s", unzippedFolderPath)
		return err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			if file.Name() == "BepInEx" {
				destination := "."
				err = pathUtil.MoveDir(filepath.Join(unzippedFolderPath, "BepInEx"), destination)
			} else if file.Name() == "plugins" {
				destination := filepath.Join(".", "BepInEx")
				err = pathUtil.MoveDir(filepath.Join(unzippedFolderPath, "plugins"), destination)
			} else if file.Name() == "patches" {
				destination := filepath.Join(".", "BepInEx")
				err = pathUtil.MoveDir(filepath.Join(unzippedFolderPath, "patchers"), destination)
			} else if file.Name() == "core" {
				destination := filepath.Join(".", "BepInEx")
				err = pathUtil.MoveDir(filepath.Join(unzippedFolderPath, "core"), destination)
			} else if file.Name() == "test" {
				destination := filepath.Join(".", "BepInEx")
				err = pathUtil.MoveDir(filepath.Join(unzippedFolderPath, "test"), destination)
			}
			if err != nil {
				return err
			}
		} else {
			destination := filepath.Join(".", "BepInEx", "plugins", file.Name())
			err = pathUtil.MoveFile(filepath.Join(unzippedFolderPath, file.Name()), destination)
			if err != nil {
				return err
			}
		}
	}

	fileInfo, err = os.ReadDir(unzippedFolderPath)
	if err != nil {
		fmt.Printf("error reading leftover mod files: %s", unzippedFolderPath)
	}
	if len(fileInfo) != 0 {
		c := color.New(color.FgHiMagenta)
		_, err := c.Printf("\nFiles leftover in dir: %s", unzippedFolderPath)
		if err != nil {
			return err
		}
		return nil
	}
	err = os.Remove(unzippedFolderPath)
	if err != nil {
		return fmt.Errorf("failed to remove source directory: %w", err)
	}
	return nil
}

func unzipMod(zipFilePath string) (string, error) {
	// zipFilePath: user/downloads/LC_New_Mods/zip/foo.zip
	zipDirectory := filepath.Dir(zipFilePath)    // user/downloads/LC_New_Mods/zip/
	unzipDirectory := filepath.Dir(zipDirectory) // user/downloads/LC_New_Mods/

	unzippedFolderPath, err := zipUtil.UnzipFile(zipFilePath, unzipDirectory)
	if err != nil {
		return "", err
	}
	return unzippedFolderPath, nil
}
