package clogs

type FileItem struct {
	Path     string // relative to domain
	Size     int64
	Checksum string
}

type CmdKind int

const (
	CmdNew CmdKind = iota
	CmdDel
)

type CmdLog struct {
	Kind      CmdKind
	FileItem  FileItem
	Timestamp int64
}

type CmdLogs struct {
	Cmds      []CmdLog
	VersionID int // same as len(Cmds)
}
