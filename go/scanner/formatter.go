package scanner

import (
	"fmt"
	"log"
	"os"
	syspath "path"
	"path/filepath"

	"github.com/xiejw/lunar/base"
)

type FileMetadata struct {
	BaseDir string      // The Walk base directory.
	Path    string      // The path (without base directory or leading '/')
	Info    os.FileInfo // `FileInfo` for the file.
}

// A function allowing to have side-effect, typically printing, operating on the
// `metadata`.
type Formatter func(metadata FileMetadata)

// Stubbed for testing.
var printf = fmt.Printf

type FormatterOption struct {
	PrintChecksum          bool
	ConvertPathToUnixSlash bool
}

func NewFormatter(option FormatterOption) Formatter {
	printChecksum := option.PrintChecksum
	pathToUnix := option.ConvertPathToUnixSlash

	return func(metadata FileMetadata) {
		info := metadata.Info
		if info.IsDir() {
			return
		}

		path := metadata.Path
		if pathToUnix {
			path = filepath.ToSlash(path)
		}
		if !printChecksum {
			printf("%12d %v\n", info.Size(), path)
			return
		}

		hash, err := base.Sha256Sum(syspath.Join(metadata.BaseDir, metadata.Path))
		if err != nil {
			log.Fatal(err)
		}
		printf("%x %v\n", hash, path)
	}
}

// Returns a `Formatter` which prints the file (not folder) and file size
// (only).
func NewPrintPathAndSizeFormatter(pathToUnix bool) Formatter {
	return NewFormatter(FormatterOption{
		PrintChecksum:          false,
		ConvertPathToUnixSlash: pathToUnix,
	})
}

// Returns a `Formatter` which prints the path and Sha256 Hash.
func NewPrintPathAndCheckSumFormatter(pathToUnix bool) Formatter {
	return NewFormatter(FormatterOption{
		PrintChecksum:          true,
		ConvertPathToUnixSlash: pathToUnix,
	})
}
