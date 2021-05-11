package main

import (
	"fmt"

	"github.com/xiejw/fsx/src/fs"
)

// Handles a single domain (filetree).
func main() {
	fmt.Printf("Hello Ft.\n")

	rootDir := "."

	// Prints a local file tree.
	ft, err := fs.FromLocalFS(rootDir)
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}

	fmt.Printf("FT: %v [\n", ft.BaseDir)
	for _, it := range ft.Items {
		fmt.Printf("%10d %v\n", it.Size, it.Path)
	}
	fmt.Printf("]\n")
}
