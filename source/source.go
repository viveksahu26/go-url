package source

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"

	"github.com/go-git/go-billy/v5"
)

func IsGit(in string) bool {
	return regexp.MustCompile("^(http|https)://").MatchString(in)
}

func IsJson(file fs.FileInfo) bool {
	if file.IsDir() {
		return false
	}
	ext := filepath.Ext(file.Name())
	fmt.Println("SBOM file extension: ", ext)
	return ext == ".json" || ext == ".spdx.json" || ext == ".spdx"
}

func ListFiles(fs billy.Filesystem, path string, predicate func(fs.FileInfo) bool) ([]string, error) {
	path = filepath.Clean(path)

	if _, err := fs.Stat(path); err != nil {
		log.Fatalf("Failed to retrieve file Info: %s", err)
	}
	// fmt.Println("Path: ", path)

	files, err := fs.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed to read path: %s", err)
	}
	// fmt.Println("files: ", files)

	var results []string

	for _, file := range files {
		name := filepath.Join(path, file.Name())

		if file.IsDir() {
			children, err := ListFiles(fs, name, predicate)
			if err != nil {
				return nil, err
			}
			results = append(results, children...)
		} else if predicate(file) {
			results = append(results, name)
		}
	}
	return results, nil
}

func ProcessPath(fs billy.Filesystem, path string) ([]string, error) {
	fileInfo, err := fs.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file info: %w", err)
	}

	var results []string
	if fileInfo.IsDir() {
		files, err := ListFiles(fs, path, IsJson)
		if err != nil {
			return nil, fmt.Errorf("failed to list files in directory: %w", err)
		}
		results = append(results, files...)
	} else {
		results = append(results, path)
	}
	return results, nil
}
