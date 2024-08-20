package zipUtil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// zipFiles zips the specified files into a single zip archive
func zipFiles(zipFileName string, files []string) error {
	// Create a new zip file
	zipfile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	// Create a new zip archive.
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// Add each file to the zip archive
	for _, file := range files {
		if err := addToZip(archive, file); err != nil {
			return err
		}
	}

	return nil
}

// addToZip adds a file or directory to the zip archive
func addToZip(archive *zip.Writer, filename string) error {
	// Open the file or directory
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		// Add the directory itself and recursively add its contents
		return addDirToZip(archive, filename)
	}

	// Create a header for the file
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Change the header name to the relative path
	header.Name = filepath.ToSlash(filename)

	// Set compression method
	header.Method = zip.Deflate

	// Create a file in the zip archive using the header
	writer, err := archive.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy the content of the file into the zip file
	if _, err = io.Copy(writer, fileToZip); err != nil {
		return err
	}

	return nil
}

// addDirToZip adds a directory and its contents to the zip archive
func addDirToZip(archive *zip.Writer, dirPath string) error {
	err := filepath.WalkDir(dirPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute the relative path of the file
		relPath, err := filepath.Rel(filepath.Dir(dirPath), path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Add directory with a trailing slash
			header := &zip.FileHeader{
				Name:   filepath.ToSlash(relPath) + "/",
				Method: zip.Store,
			}
			_, err := archive.CreateHeader(header)
			return err
		} else {
			// Add file
			return addToZip(archive, path)
		}
	})
	return err
}

func getZipPackVersionFromUser() (string, error) {
	// Define the regex pattern to match a number in the format # or #.#
	pattern := `^\d+(\.\d+)?$`
	re := regexp.MustCompile(pattern)

	for {
		var num string
		fmt.Printf("What version is this pack (1, 2, etc.): ")
		_, err := fmt.Scanln(&num)
		if err != nil {
			return "", fmt.Errorf("error reading: %v", err)
		}

		if re.MatchString(num) {
			return num, nil
		} else {
			fmt.Println("Invalid input. Please enter a valid version number.")
		}
	}
}

func ZipBepinEx() error {
	//Files and directories to be zipped
	files := []string{"winhttp.dll", "doorstop_config.ini", "BepInEx"}

	version, err := getZipPackVersionFromUser()
	if err != nil {
		return fmt.Errorf("could not get zip pack version from user: %e", err)
	}

	zipFileName := fmt.Sprintf("BepinExPack_v%s.zip", version)

	// Create zip file
	err = zipFiles(zipFileName, files)
	if err != nil {
		return fmt.Errorf("error zipping files: %e", err)
	}
	return nil
}
