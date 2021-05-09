package main

import (
	"fmt"

	"github.com/xiejw/fsx/src/fs"
)

func main() {
	fmt.Printf("Hello FsX.\n")

	// Prints a local file tree.
	ft, err := fs.FromLocalFS(".")
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
