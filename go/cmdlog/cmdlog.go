package cmdlog

type CmdLog struct {
	Version    uint64
	MasterName string
	Timestamp  int64
	Cmd        *Cmd
}

type Cmd struct {
	Type     CmdType
	Dir      string // Dir == "" -> root
	FileName string // Required.
	Size     uint64 // Must be filled for CmdNew
	Checksum []byte // Must be filled for CmdNew
}

type CmdType int

const (
	CmdUnspecified CmdType = iota
	CmdNew
	CmdDelete
)
