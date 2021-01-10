package snapshot

type (
	// Represents a group of FileItems forming a file tree. Typically it could be a real file tree on
	// the disk.
	Region struct {
		Name    string     // Symbol name for the region.
		BaseDir string     // The Base dir. Root for all Items.
		Items   []FileItem // All items in the region. Guaranteed to be sorted alphabetically.
	}

	FileItem struct {
		FullPath string // Full path, but relative to Base dir.
		Size     uint64 // Length in bytes for regular files.
		Checksum []byte // Optional sha256 checksum. See Snapshot interface for requirement.
	}
)
