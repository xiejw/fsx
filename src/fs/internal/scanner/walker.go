// Package scanner scans a folder and allows to take action on each item it finds.
package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xiejw/fsx/src/errors"
)

// -------------------------------------------------------------------------------------------------
// public
// -------------------------------------------------------------------------------------------------

// The metadata emitted by Walk.
type FileMetadata struct {
	BaseDir string      // The base directory.
	Path    string      // The relative path (without base directory or leading '/')
	Info    os.FileInfo // `FileInfo` for the file.
}

// A function allowing to have side-effect, typically printing, operating on the `metadata`.
type Formatter func(metadata FileMetadata)

// Filter is invoked for each entry during the tree walk. Returns true to skip the
// item passing to formatter.
//
// If any of the `filters` returns true for the sub-folder, the entire sub-tree
// is skipped.
type Filter func(path string, info os.FileInfo) bool

// Walk walks the file tree recursively rooted at `baseDir` in lexical order.
//
// - If none of the `filters` returns true, the file item, including folder, will be passed to
//   `formatter`.
// - If any of the `filters` returns true for the sub-folder, the entire sub-tree is skipped.
func Walk(baseDir string, filters []Filter, formatter Formatter) error {

	// Use absolute path to avoid starting a baseDir == `.`, which is considered as hidden file.
	dir, err := filepath.Abs(baseDir)
	if err != nil {
		return err
	}

	// `Walk` does not support follow link. So, we read the content of the link if possible.
	realDirPath := mustFollowDirLink(dir)
	return walkImpl(realDirPath, filters, formatter)
}

// -------------------------------------------------------------------------------------------------
// impl
// -------------------------------------------------------------------------------------------------

// stubbed for testing.
var filePathWalk = filepath.Walk

func walkImpl(baseDir string, filters []Filter, formatter Formatter) error {

	// Defines a walkFn which captured information.
	walkFn := func(path string, info os.FileInfo, err error) error {
		shouldSkip := false
		for _, filter := range filters {
			if filter(path, info) {
				if info.IsDir() {
					// Skips the folder.
					return filepath.SkipDir
				} else {
					// Skips the current file but continues in the folder.
					shouldSkip = true
					break
				}
			}
		}
		if !shouldSkip {
			metadata := FileMetadata{
				BaseDir: baseDir,
				Path:    strings.TrimPrefix(strings.TrimPrefix(path, baseDir), "/"),
				Info:    info,
			}
			formatter(metadata)
		}
		return nil
	}

	if err := filePathWalk(baseDir, walkFn); err != nil {
		return errors.WrapNote(err, "failed to walk into dir: %s", baseDir)
	}
	return nil
}

func mustFollowDirLink(dir string) string {
	stat, err := os.Lstat(dir)
	if err != nil {
		panic(fmt.Sprintf("failed to stat the dir (%v): %v", dir, err))
	}

	if stat.Mode()&os.ModeSymlink != 0 {
		realDir, err := os.Readlink(dir)
		if err != nil {
			panic(fmt.Sprintf("failed to read link for dir %v: %v", dir, err))
		}
		return realDir
	}
	return dir
}
