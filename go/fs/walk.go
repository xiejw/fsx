package fs

import (
	"github.com/xiejw/lunar/scanner"
)

type WalkResult struct {
	Path string
	Size uint64
	Err  error
}

// Walks the tree at `baseDir` and performs `cb` on each file item.
func Walk(baseDir string, cb func(*WalkResult)) error {
	err := scanner.Walk(baseDir,
		[]scanner.Filter{scanner.NewCommonFilter([]string{})},
		func(metadata scanner.FileMetadata) {
			info := metadata.Info
			if info.IsDir() {
				return
			}
			cb(&WalkResult{
				Path: metadata.Path,
				Size: uint64(info.Size()),
				Err:  nil, // No error for local fs.
			})
		})
	return err
}
