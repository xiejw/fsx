// Package scanner scans a folder and allows to take action on each item it finds.
package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

// Stubbed for testing.
var filePathWalk = filepath.Walk

// Walks the file tree rooted at `baseDir`. If none of the `filters` returns
// true, the file item, including folder, will be passed to `formatter`.
//
// If any of the `filters` returns true for the sub-folder, the entire sub-tree
// is skipped.
//
// The files/folders are walked in lexical order.
func Walk(baseDir string, filters []Filter, formatter Formatter) {
	// Walk does not support follow link. So, we read the content of the link if possible.
	realDirPath := mustFollowDirLink(baseDir)
	walkImpl(realDirPath, filters, formatter)
}

func walkImpl(baseDir string, filters []Filter, formatter Formatter) {

	// Defines a walkFn wich captured information.
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
		glog.Fatalf("Failed to walk: %v", err)
	}
}

func mustFollowDirLink(dir string) string {
	stat, err := os.Lstat(dir)
	if err != nil {
		glog.Fatalf("failed to stat the dir (%v): %v", dir, err)
		return dir
	}

	if stat.Mode()&os.ModeSymlink != 0 {
		realDir, err := os.Readlink(dir)
		if err != nil {
			glog.Fatalf("failed to read link for dir %v: %v", dir, err)
			return dir
		}
		glog.Infof("read link %v => %v", dir, realDir)
		return realDir
	}
	return dir
}
