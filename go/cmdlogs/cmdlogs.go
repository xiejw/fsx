package cmdlogs

import (
	"io"
)

type NodeState struct {
	Name        string   // Must be uinque in the cluster.
	NextVersion uint64   // Points to next verison.
	CmdLogs     []CmdLog // Ordered CmdLog
}

type CmdType int

const (
	CmdUnspecified CmdType = iota
	CmdNew
	CmdDelete
)

type Cmd struct {
	Type     CmdType
	Dir      string // Dir == "" -> root
	FileName string // Must be file. Cannot be empty.
	Checksum []byte // Must be filled for CmdNew
}

type CmdLog struct {
	Version    uint64
	MasterName string
	Cmd        Cmd
}

func (state *NodeState) Markshal(w io.Writer) err {
	return nil
}
