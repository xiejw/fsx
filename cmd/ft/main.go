package main

import (
	"fmt"

	"github.com/xiejw/fsx/src/errors"
	"github.com/xiejw/fsx/src/fs"
)

// Handles a single domain (filetree).
func main() {
	fmt.Printf("Hello Ft.\n")

	rootDir := "."
	_, err := fetchLocalFS(rootDir)
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}
}

// -------------------------------------------------------------------------------------------------
// impl
// -------------------------------------------------------------------------------------------------
func fetchLocalFS(rootDir string) (*fs.FileTree, error) {
	// Prints a local file tree.
	ft, err := fs.FromLocalFS(rootDir)
	if err != nil {
		return nil, errors.WrapNote(err, "failed to fetch local file tree at: %v", rootDir)
	}

	maxItems := 10

	fmt.Printf("FT: %v [\n", ft.BaseDir)
	for i, it := range ft.Items {
		fmt.Printf("%10d %v\n", it.Size, it.Path)
		if i == maxItems-1 && i != len(ft.Items)-1 {
			fmt.Printf("       ... (%v items left)\n", len(ft.Items)-maxItems)
			break
		}
	}
	fmt.Printf("] (%v total items)\n", len(ft.Items))

	return ft, nil
}
