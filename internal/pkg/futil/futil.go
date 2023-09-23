package futil

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func DirHasFile(rootDirPath string, filename string) bool {
	file, _ := FindFileAbsPath(rootDirPath, filename)
	return file != ""
}

func FindFileAbsPath(rootDirPath string, filename string) (string, error) {
	file := ""

	filepath.Walk(rootDirPath, func(path string, fileInfo os.FileInfo, err error) error {
		if file != "" {
			return filepath.SkipAll
		}

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
	}

	return file, nil
}

func FindFileSha256(rootDirPath string, filename string) (string, error) {
	fileAbsPath, err := FindFileAbsPath(rootDirPath, filename)

	if err != nil {
		return "", err
	}

	return FileSha256(fileAbsPath)
}

func FileSha256(fileAbsPath string) (string, error) {
	file, err := os.Open(fileAbsPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
