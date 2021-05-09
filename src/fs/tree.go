package fs

import (
	"github.com/xiejw/fsx/src/errors"
	"github.com/xiejw/fsx/src/fs/scanner"
)

// same as clogs.FileItem
type FileItem struct {
	Path     string // relative to domain.
	Size     int64  // file size.
	Checksum string // optionel checksum.
}

type FileTree struct {
	BaseDir     string      // domain path.
	HasChecksum bool        // if false, all items do not have checksum.
	Items       []*FileItem // alphabetically sorted items.
}

// creates a FileTree by walking the baseDir.
func FromLocalFS(baseDir string) (*FileTree, error) {
	ft := &FileTree{
		BaseDir:     baseDir,
		HasChecksum: false, // no checksum for perf.
		Items:       nil,
	}

	items := make([]*FileItem, 0)

	cb := func(metadata scanner.FileMetadata) {
		info := metadata.Info
		if info.IsDir() {
			return
		}
		items = append(items, &FileItem{
			Path: metadata.Path,
			Size: info.Size(),
		})
	}

	err := scanner.Walk(baseDir, []scanner.Filter{scanner.NewCommonFilter(nil)}, cb)
	if err != nil {
		return nil, errors.WrapNote(err, "failed to create FileTree by walking local fs.")
	}

	ft.Items = items

	return ft, nil
}
