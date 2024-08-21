package pathUtil

import (
	"fmt"
	"io"
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

// MoveDir moves a directory and its contents from dirToBeMoved to destinationDirParent.
//
// Parameters:
//   - dirToBeMoved: The path of the directory to be moved. This can be an absolute
//     or relative path.
//   - destinationDirParent: The parent directory where the source directory should
//     be moved. The source directory will be placed inside this parent directory
//     with the same name.
func MoveDir(dirToBeMoved string, destinationDirParent string, dirNewName ...string) error {
	// Get properties of the source directory
	srcInfo, err := os.Stat(dirToBeMoved)
	if err != nil {
		return fmt.Errorf("failed to get properties of source directory: %w", err)
	}

	var destDir string
	if len(dirNewName) == 0 {
		// Create the destination directory with the same name as the source directory
		destDir = filepath.Join(destinationDirParent, filepath.Base(dirToBeMoved))
	} else if len(dirNewName) == 1 {
		destDir = filepath.Join(destinationDirParent, filepath.Base(dirNewName[0]))
	} else {
		return fmt.Errorf("too many new name arguments")
	}
	err = os.MkdirAll(destDir, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Iterate over the files and directories in the source directory
	err = filepath.Walk(dirToBeMoved, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Compute the relative path from the source directory
		relativePath, err := filepath.Rel(dirToBeMoved, srcPath)
		if err != nil {
			return fmt.Errorf("failed to compute relative path: %w", err)
		}

		// Compute the destination path
		dstPath := filepath.Join(destDir, relativePath)

		if info.IsDir() {
			// Create directories in the destination
			err := os.MkdirAll(dstPath, info.Mode())
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dstPath, err)
			}
		} else {
			// Move files
			err := MoveFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path: %w", err)
	}

	// Optionally, remove the now-empty source directory
	err = os.RemoveAll(dirToBeMoved)
	if err != nil {
		return fmt.Errorf("failed to remove source directory: %w", err)
	}

	return nil
}

// MoveFile moves a file from the specified source file path to the destination path.
//
// Parameters:
//   - sourceFilePath: The path of the file to be moved. This can be a local file name
//     or a full file path.
//   - destinationPath: The destination where the file should be moved. This can be
//     either a directory path or a full file path. If a directory path is provided,
//     the file will be moved into this directory with the same name as the source file.
func MoveFile(sourceFilePath string, destinationFilePath string) error {
	// Copy the file to the destination
	err := copyFile(sourceFilePath, destinationFilePath)
	if err != nil {
		return err
	}
	// Remove the original file
	err = os.Remove(sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	return nil
}

// copyFile copies a file from the specified source file path to the destination path.
//
// Parameters:
//   - sourceFilePath: The path of the source file to be copied. This can be a local
//     file name or a full file path.
//   - destinationPath: The destination where the file should be copied. This can
//     either be a directory path or a full file path. If a directory path is provided,
//     the file will be copied into this directory with the same name as the source file.
func copyFile(sourceFilePath string, destinationPath string) error {
	// Get the file name from the source file path
	fileName := filepath.Base(sourceFilePath)

	// Check if the destination path is a directory
	destInfo, err := os.Stat(destinationPath)
	if err == nil && destInfo.IsDir() {
		destinationPath = filepath.Join(destinationPath, fileName)
	}

	srcFile, err := os.Open(sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Copy file permissions
	srcInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}
	err = os.Chmod(destinationPath, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}
