package ffind

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindFileAbsPath(rootDirPath string, filename string) (string, error) {
	file := ""

	filepath.Walk(rootDirPath, func(path string, fileInfo os.FileInfo, err error) error {
		fileInfo, err = os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.Mode().IsRegular() && fileInfo.Name() == filename {
			file = path
		}

		return nil
	})

	if file == "" {
		return file, fmt.Errorf("File not found: %s", filename)
	} else {
		return file, nil
	}
}
