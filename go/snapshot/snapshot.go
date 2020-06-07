package snapshot

import "errors"

var (
	ErrItemNotExist     = errors.New("item not exist.")
	ErrItemAlreadyExist = errors.New("item already exist.")
)

type FileItem struct {
	FullPath string // Full path, but relative to base dir.
	Checksum []byte // Allows to be empty.
}

// Order is not guaranteed.
type Iterator interface {
	Next() *FileItem
}

// Represents a snaptshot of the current system.
type Snapshot interface {
	LookUp(fullPath string) *FileItem // Returns nil for not exist.
	Add(*FileItem) error
	Delete(fullPath string) error
	Iterator() Iterator
}
