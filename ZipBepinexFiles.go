package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

	// Create a header for the file or directory
	var header *zip.FileHeader
	if info.IsDir() {
		// If it's a directory, use FileInfoHeader to create the header
		header, err = zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
	} else {
		// If it's a file, also use FileInfoHeader
		header, err = zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
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

	// If it's a file, copy its content into the zip file
	if !info.IsDir() {
		if _, err = io.Copy(writer, fileToZip); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	//Files and directories to be zipped
	files := []string{"winhttp.dll", "doorstop_config.ini", "BepInEx"}

	var num string
	fmt.Printf("What version is this pack (1, 2, etc.): ")
	_, err := fmt.Scanln(&num)
	if err != nil {
		fmt.Println("error read")
		quit()
	}
	version, err := strconv.Atoi(num)
	if err != nil {
		fmt.Println("not an int")
		quit()
	}

	zipFileName := fmt.Sprintf("BepinExPack_v%d.zip", version)
	fmt.Println(zipFileName)

	// Create zip file
	err = zipFiles(zipFileName, files)
	if err != nil {
		fmt.Println("Error zipping files:", err)
		quit()
	}

	fmt.Println("Files zipped successfully.")
	quit()
}

func quit() {
	fmt.Printf("Press 'enter' to exit...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
	os.Exit(0)
}
