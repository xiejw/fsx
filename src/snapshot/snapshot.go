package snapshot

import (
	"github.com/xiejw/fsx/src/errors"
)

type (
	// Represents a group of FileItems forming a file tree. Typically it could be a real file tree on
	// the disk.
	Region struct {
		Name    string      // Symbol name for the region.
		BaseDir string      // The Base dir. Root for all Items.
		Items   []*FileItem // All items in the region. Guaranteed to be sorted alphabetically.
	}

	FileItem struct {
		RelPath  string // Full path, but relative to Base dir.
		Size     uint64 // Length in bytes for regular files.
		Checksum []byte // Optional sha256 checksum. See Snapshot interface for requirement.
	}
)

func FillRegion(r *Region) error {
	if r.BaseDir == "" {
		return errors.New("BaseDir cannot be empty to fill region")
	}
	if len(r.Items) != 0 {
		return errors.New("Items must be empty to fill region")
	}

	cb := func(re *walkResult) {
		r.Items = append(r.Items, &FileItem{
			RelPath:  re.Path,
			Size:     re.Size,
			Checksum: nil,
		})
	}

	return walk(r.BaseDir, cb)
}
