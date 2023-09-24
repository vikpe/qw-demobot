package futil

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FindFilepathsByExtension(rootDirPath string, ext string) []string {
	files := []string{}

	filepath.Walk(rootDirPath, func(path string, fileInfo os.FileInfo, err error) error {
		fileInfo, err = os.Stat(path)
		if err != nil {
			return err
		}

		if fileInfo.Mode().IsRegular() && strings.HasSuffix(path, ext) {
			files = append(files, path)
		}

		return nil
	})

	return files
}

func DirHasFile(rootDirPath string, filename string) bool {
	file, _ := FindFilepath(rootDirPath, filename)
	return file != ""
}

func FindFilepath(rootDirPath string, filename string) (string, error) {
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
	fileAbsPath, err := FindFilepath(rootDirPath, filename)

	if err != nil {
		return "", err
	}

	return FileSha256(fileAbsPath)
}

func FileSha256(filepath string) (string, error) {
	file, err := os.Open(filepath)
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
