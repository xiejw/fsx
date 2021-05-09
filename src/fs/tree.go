package fs

// same as clogs.FileItem
type FileItem struct {
	Path     string // relative to domain.
	Size     int64  // file size.
	Checksum string // optionel checksum.
}

type FileTree struct {
	BaseDir     string     // domain path.
	HasChecksum bool       // if false, all items do not have checksum.
	Items       []FileItem // alphabetically sorted items.
}
