package snapshot

import (
	"github.com/xiejw/fsx/src/fs/scanner"
)

type walkResult struct {
	Path string
	Size uint64
	Err  error
}

// Walks the tree at `baseDir` and performs `cb` on each file item.
func walk(baseDir string, cb func(*walkResult)) error {
	err := scanner.Walk(baseDir,
		[]scanner.Filter{scanner.NewCommonFilter([]string{})},
		func(metadata scanner.FileMetadata) {
			info := metadata.Info
			if info.IsDir() {
				return
			}
			cb(&walkResult{
				Path: metadata.Path,
				Size: uint64(info.Size()),
				Err:  nil, // No error for local fs.
			})
		})
	return err
}
