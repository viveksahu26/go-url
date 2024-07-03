package main

import (
	"fmt"

	"github.com/viveksahu26/go-url/source"
)

func main() {
	url := "https://github.com/interlynk-io/sbomqs/blob/main/samples/sbomqs-spdx-syft.json"

	if source.IsGit(url) {
		fmt.Println("Yes, it's a git url: ", url)
	}
}
