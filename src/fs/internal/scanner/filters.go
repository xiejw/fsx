package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// Returns true for the following patterns:
//
//   - hidden folder or files..
//   - *.swp file,
func NewCommonFilter(extsToSkip []string) Filter {
	// adds vim swp ext.
	extsToSkip = append(extsToSkip, "swp")
	dotExtsToSkip := make([]string, 0, len(extsToSkip))
	for _, ext := range extsToSkip {
		if ext == "" {
			continue
		}
		dotExtsToSkip = append(dotExtsToSkip, "."+ext)
	}

	return func(path string, info os.FileInfo) bool {
		_, file := filepath.Split(path)
		if strings.HasPrefix(file, ".") {
			// Skips for hidden directory and file.
			// glog.V(1).Infof("skip %v as it is hidden.", path)
			return true
		}

		ext := filepath.Ext(path)
		for _, ext_to_skip := range dotExtsToSkip {
			if ext == ext_to_skip {
				// glog.V(1).Infof("skip %v as it has extension (%v).", path, ext)
				return true
			}
		}
		return false
	}
}

// Returns true if number of items srreaming into this Filter is not larger
// than the `total_num`. If `fileOnly`, only counts the files.
func NewCounterFilter(total_num int, fileOnly bool) Filter {
	count := 0
	counterFilter := func(path string, info os.FileInfo) bool {
		if fileOnly {
			if !info.IsDir() {
				count += 1
			}
		} else {
			count += 1
		}
		return count > total_num
	}
	return counterFilter
}
