package main

import (
	"fmt"
	"os"
	"path"

	"github.com/xiejw/fsx/src/clogs"
	"github.com/xiejw/fsx/src/errors"
	"github.com/xiejw/fsx/src/fs"
)

var flagClogsFile = "clogs.txt"

// Handles a single domain (filetree).
func main() {
	fmt.Printf("Hello Ft.\n")

	rootDir := "."
	ft_local, err := fetchFtFromLocalDir(rootDir)
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}
	printFileTree("local", ft_local)
	Printf("\n")

	ft_clgs, err := fetchFtFromClogsFile(rootDir, flagClogsFile)
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}
	printFileTree("clogs", ft_clgs)
}

// -------------------------------------------------------------------------------------------------
// impl
// -------------------------------------------------------------------------------------------------
func fetchFtFromLocalDir(rootDir string) (*fs.FileTree, error) {
	// Prints a local file tree.
	ft, err := fs.FromLocalFS(rootDir)
	if err != nil {
		return nil, errors.WrapNote(err, "failed to fetch local file tree at: %v", rootDir)
	}

	return ft, nil
}

func fetchFtFromClogsFile(rootDir, clogsPath string) (*fs.FileTree, error) {
	var clgs *clogs.CmdLogs

	absPath := path.Join(rootDir, clogsPath)
	_, err := os.Stat(absPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.WrapNote(
				err, "failed to create file tree by loading clogs file: %s", clogsPath)
		}

		// empty clgs if fine.
		clgs = &clogs.CmdLogs{}
	} else {
		f, err := os.Open(absPath)
		if err != nil {
			return nil, errors.WrapNote(
				err, "failed to create file tree by loading clogs file: %s", clogsPath)
		}
		defer f.Close()

		clgs, err = clogs.FromLines(f)
		if err != nil {
			return nil, errors.WrapNote(
				err, "failed to create file tree by loading clogs file: %s", clogsPath)
		}
	}

	ft, err := fs.FromCmdLogs(rootDir, clgs)
	if err != nil {
		return nil, errors.WrapNote(
			err, "failed to create file tree by loading clogs file: %s", clogsPath)
	}

	return ft, nil
}

func printFileTree(startingMsg string, ft *fs.FileTree) {
	maxItems := 10

	fmt.Printf("FT (%v): %v [\n", startingMsg, ft.BaseDir)
	for i, it := range ft.Items {
		fmt.Printf("%10d %v\n", it.Size, it.Path)
		if i == maxItems-1 && i != len(ft.Items)-1 {
			fmt.Printf("       ... (%v items left)\n", len(ft.Items)-maxItems)
			break
		}
	}
	fmt.Printf("] (%v total items)\n", len(ft.Items))
}
