package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/viveksahu26/go-url/source"
)

func main() {
	url_path := "https://github.com/interlynk-io/sbomqs/blob/main/samples/sbomqs-spdx-syft.json"

	if source.IsGit(url_path) {
		fmt.Println("Yes, it's a git url: ", url_path)

		fs := memfs.New()
		gitURL, err := url.Parse(url_path)
		if err != nil {
			log.Fatalf("err: ", err)
		} else {
			fmt.Println("parse gitURL: ", gitURL)
			pathElems := strings.Split(gitURL.Path[1:], "/")
			fmt.Println("pathElems: ", pathElems)

			if len(pathElems) <= 1 {
				log.Fatalf("invalid URL path %s - expected https://github.com/:owner/:repository/:branch (without --git-branch flag) OR https://github.com/:owner/:repository/:directory (with --git-branch flag)", gitURL.Path)
			}

			fmt.Println("gitURL.Path: ", gitURL.Path)

			gitURL.Path = strings.Join([]string{pathElems[0], pathElems[1]}, "/")
			fmt.Println("gitURL.Path: ", gitURL.Path)

			repoURL := gitURL.String()
			fmt.Println("repoURL: ", repoURL)

			filePath := strings.Join([]string{pathElems[len(pathElems)-2], pathElems[len(pathElems)-1]}, "/")
			fmt.Println("filePath: ", filePath)

			cloneOptions := &git.CloneOptions{
				URL:           repoURL,
				ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", "main")),
				Depth:         1,
				Progress:      os.Stdout,
				SingleBranch:  true,
			}

			repo, err := git.Clone(memory.NewStorage(), fs, cloneOptions)
			if err != nil {
				log.Fatalf("Failed to clone repository: %s", err)
			}
			fmt.Println("repo: ", repo)

			gitPathToJson := "/"

			allFiles, err := source.ListFiles(fs, gitPathToJson, source.IsJson)
			if err != nil {
				log.Fatalf("error: failed to list all files in repository (%w)", err)
				fmt.Println("error: failed to list all files in repository (%w)", err)
			}
			fmt.Println("all Files: ", allFiles)
			var ourFiles []string

			for _, file := range allFiles {
				if strings.Contains(file, filePath) {
					fmt.Println("Yes, this file is ours")
					ourFiles = append(ourFiles, file)
				}
			}
			fmt.Println("ourFiles: ", ourFiles)

			for _, file := range ourFiles {
				filePath, err := fs.Open(file)
				if err != nil {
					fmt.Println("error: failed to open file(%w)", err)
				}
				data, err := io.ReadAll(filePath)
				if err != nil {
					fmt.Println("error: failed to read file(%w)", err)
				}
				fmt.Println("data: ", string(data))

			}

		}

	}
}
