package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/xiejw/lunar/base/errors"

	"github.com/xiejw/fsx/src/clogs"
	"github.com/xiejw/fsx/src/fs"
)

var (
	flagClogsFile = "clogs.txt"
	flagChecksum  = true
)

type Config struct {
	RootDir               string
	LoadChecksum          bool
	ClogsFile             string
	DiffFS                bool
	GenerateCmdLogsChange bool // Need DiffFS
	SyncToFile            bool // Need GenerateCmdLogsChange
	PrintLocalFS          bool
	PrintClogsFS          bool
	PrintDiff             bool
	PrintCmdLogsChange    bool
}

// Handles a single domain (filetree).
func main() {
	fmt.Printf("Hello Ft.\n")

	config := Config{
		RootDir:               ".",
		LoadChecksum:          flagChecksum,
		ClogsFile:             flagClogsFile,
		DiffFS:                true,
		GenerateCmdLogsChange: true,
		SyncToFile:            true,
		PrintLocalFS:          true,
		PrintClogsFS:          true,
		PrintDiff:             true,
		PrintCmdLogsChange:    true,
	}

	ft_local, err := fetchFtFromLocalDir(config.RootDir, config.LoadChecksum)
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}
	if config.PrintLocalFS {
		printFileTree("local", ft_local)
		fmt.Printf("\n")
	}

	ft_clgs, err := fetchFtFromClogsFile(config.RootDir, config.ClogsFile)
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}
	if config.PrintClogsFS {
		printFileTree("clogs", ft_clgs)
	}

	if !config.DiffFS {
		return
	}

	diffResult, err := ft_clgs.ConvertTo(ft_local)
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}

	if config.PrintDiff {
		diffResult.Print(os.Stdout)
	}

	if !config.GenerateCmdLogsChange {
		return
	}

	clgs, err := diffResult.ToCmdLogs(time.Now().Unix())
	if err != nil {
		fmt.Printf("unexpected error: %v", err)
		return
	}

	if config.PrintCmdLogsChange {
		fmt.Printf("CLogs changes:\n")
		if len(clgs) == 0 {
			fmt.Printf("(empty)\n")
		}
		for _, clg := range clgs {
			fmt.Printf("%+v\n", clg.ToOneLine())
		}
	}
}

// -------------------------------------------------------------------------------------------------
// impl
// -------------------------------------------------------------------------------------------------
func fetchFtFromLocalDir(rootDir string, checksum bool) (*fs.FileTree, error) {
	// Prints a local file tree.
	ft, err := fs.FromLocalFS(rootDir, checksum)
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
	checksum := ft.HasChecksum

	for i, it := range ft.Items {
		if checksum {
			fmt.Printf("%10d %v %v\n", it.Size, it.Checksum, it.Path)
		} else {
			fmt.Printf("%10d %v\n", it.Size, it.Path)
		}
		if i == maxItems-1 && i != len(ft.Items)-1 {
			fmt.Printf("       ... (%v items left)\n", len(ft.Items)-maxItems)
			break
		}
	}
	fmt.Printf("] (%v total items)\n", len(ft.Items))
}
