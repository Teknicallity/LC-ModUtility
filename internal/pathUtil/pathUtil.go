package pathUtil

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetDownloadFolderPath() string {
	var fullDownloadDirectory string
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory\n")
	}

	downloadDirNames := []string{"Downloads", "downloads", "download", "Download"}

	for _, ddn := range downloadDirNames {
		downloadDirectory := filepath.Join(homeDirectory, ddn)

		if _, err := os.Stat(downloadDirectory); os.IsNotExist(err) {
			fullDownloadDirectory = ""
		} else {
			fullDownloadDirectory = downloadDirectory
			break
		}
	}

	return fullDownloadDirectory
}
