package snapshot

import "errors"

var (
	ErrItemNotExist     = errors.New("item not exist.")
	ErrItemAlreadyExist = errors.New("item already exist.")
)

type FileItem struct {
	FullPath string // Full path, but relative to base dir.
	Size     uint64 // Length in bytes for regular files.
	Checksum []byte // Optional sha256 checksum. See Snapshot interface for requirement.
}

// Order is not guaranteed.
type Iterator interface {
	Next() *FileItem
}

// Represents a snaptshot of the current system.
type Snapshot interface {
	LookUp(fullPath string) *FileItem // Returns nil for not exist.
	Size() uint64 // Num of FileItem.
	Iterator() Iterator // Returns a one-off iterator.

	// All FileItems in the Snapshot must all have Checksum
	HasChecksum() bool
}

type SnapshotBuilder interface {
	Snapshot

	Add(*FileItem) error
	Delete(fullPath string) error
}
