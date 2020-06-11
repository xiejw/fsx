package cmdlog

import "fmt"

type CmdLogs []CmdLog

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
	CmdDel
)

// Performs sanitys check on the CmdLogs.
func (logs CmdLogs) Check() error {
	for i, cmdLog := range []CmdLog(logs) {
		if cmdLog.Version != uint64(i) {
			return fmt.Errorf("at position %v, the CmdLog version is not right: %v.",
				i, cmdLog.Version)
		}
		if cmdLog.Timestamp == 0 {
			return fmt.Errorf("timestamp must be set.")
		}
		if len(cmdLog.MasterName) == 0 {
			return fmt.Errorf("mastername must be set.")
		}
	}

	return nil
}
