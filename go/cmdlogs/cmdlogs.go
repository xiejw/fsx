package cmdlogs

type CmdLog struct {
	Version    uint64
	MasterName string
	Cmd        *Cmd
}

type Cmd struct {
	Type     CmdType
	Dir      string // Dir == "" -> root
	FileName string // Must be file. Cannot be empty.
	Checksum []byte // Must be filled for CmdNew
}

type CmdType int

const (
	CmdUnspecified CmdType = iota
	CmdNew
	CmdDelete
)
