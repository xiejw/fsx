package fs

// same as clogs.FileItem
type FileItem struct {
	Path     string // relative to domain
	Size     int64
	Checksum string
}

type FileTree struct {
	BaseDir string
	Items   []FileItem
}
