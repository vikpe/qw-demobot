package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Tournament struct {
	Name   string  `json:"name"`
	Events []Event `json:"events"`
}

type Event struct {
	Name   string `json:"name"`
	BestOf struct {
		t4on4Final int64 `json:"4on4_final"`
	} `json:"best_of"`
}

func findFilesByExt(dir_path string, ext string) []string {
	files := []string{}

	filepath.Walk(dir_path, func(path string, fileInfo os.FileInfo, err error) error {
		fileInfo, err = os.Stat(path)
		if err != nil {
			panic(err)
		}

		if strings.HasSuffix(path, ext) {
			if fileInfo.Mode().IsRegular() {
				relativeFilePath := strings.TrimPrefix(path, dir_path)
				files = append(files, relativeFilePath)
			}
		}

		return nil
	})

	return files
}

func main() {
	files := findFilesByExt("/home/vikpe/games/demoquake/", "mvd")

	for _, i := range files {
		fmt.Println(i)
	}
}
