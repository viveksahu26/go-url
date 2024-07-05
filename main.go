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
	// urlPath := "https://github.com/interlynk-io/sbomqs/blob/main/samples"
	urlPath := "https://github.com/chainguard-dev/bom-shelter/tree/main/in-the-wild/spdx"

	if source.IsGit(urlPath) {
		fmt.Println("Yes, it's a git url: ", urlPath)

		fs := memfs.New()

		gitURL, err := url.Parse(urlPath)
		if err != nil {
			log.Fatalf("err:%v ", err)
		}

		fmt.Println("parse gitURL: ", gitURL)

		pathElems := strings.Split(gitURL.Path[1:], "/")
		if len(pathElems) <= 1 {
			log.Fatalf("invalid URL path %s - expected https://github.com/:owner/:repository/:branch (without --git-branch flag) OR https://github.com/:owner/:repository/:directory (with --git-branch flag)", gitURL.Path)
		}

		fmt.Println("pathElems: ", pathElems)
		fmt.Println("Before gitURL.Path: ", gitURL.Path)

		gitURL.Path = strings.Join([]string{pathElems[0], pathElems[1]}, "/")
		fmt.Println("After gitURL.Path: ", gitURL.Path)

		repoURL := gitURL.String()
		fmt.Println("repoURL: ", repoURL)

		fileOrDirPath := strings.Join(pathElems[4:], "/")
		fmt.Println("lastPathElement: ", fileOrDirPath)

		cloneOptions := &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", "main")),
			Depth:         1,
			Progress:      os.Stdout,
			SingleBranch:  true,
		}

		_, err = git.Clone(memory.NewStorage(), fs, cloneOptions)
		if err != nil {
			log.Fatalf("Failed to clone repository: %s", err)
		}

		var paths []string
		if paths, err = source.ProcessPath(fs, fileOrDirPath); err != nil {
			log.Fatalf("Error processing path: %v", err)
		}

		for _, p := range paths {
			fmt.Println("File Path:", p)
			dataOpen, err := fs.Open(p)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			data, err := io.ReadAll(dataOpen)
			if err != nil {
				fmt.Println("error: failed to read file", err)
				continue
			}
			fmt.Println("Data from file:", string(data))

		}

	}
}
