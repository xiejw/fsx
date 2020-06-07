package cmdlogs

import (
	"io"
)

type NodeState struct {
	Name        string   // Must be uinque in the cluster.
	NextVersion uint64   // Points to next verison.
	CmdLogs     []CmdLog // Ordered CmdLog
}

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

func (state *NodeState) Check() error {
	return nil
}

func (state *NodeState) Marshal(w io.Writer) error {
	return nil
}
