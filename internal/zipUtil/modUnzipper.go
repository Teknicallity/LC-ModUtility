package zipUtil

import (
	"path/filepath"
)

func UnzipMods(zipFilePaths []string) error {
	for _, zipFilePath := range zipFilePaths {

		err := UnzipMod(zipFilePath)
		if err != nil {
			return err
		}

	}
	return nil
}

func UnzipMod(zipFilePath string) error {
	// zipFilePath: user/downloads/LC_New_Mods/zip/foo.zip
	zipDirectory := filepath.Dir(zipFilePath)    // user/downloads/LC_New_Mods/zip/
	unzipDirectory := filepath.Dir(zipDirectory) // user/downloads/LC_New_Mods/

	err := unzipFile(zipFilePath, unzipDirectory)
	if err != nil {
		return err
	}
	return nil
}
