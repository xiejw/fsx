package cmdlog

import "fmt"

type CmdLogs []CmdLog

type CmdLog struct {
	Version     uint64 // version of the log.
	MasterName  string // name of the master node.
	Timestamp   int64  // seconds since epoch.
	LogChecksum []byte // checksum of the cmdlog.
	Cmd         *Cmd
}

type Cmd struct {
	Type     CmdType // cannot be unspecified.
	Dir      string  // empty string ("") means root dir.
	FileName string  // required.
	Size     uint64  // required for CmdNew.
	Checksum []byte  // required for CmdNew.
}

type CmdType int

const (
	CmdUnspecified CmdType = iota
	CmdNew
	CmdDel
)

// Performs sanitys check on the CmdLogs.
func (logs CmdLogs) Check() error {
	for i, cmdLog := range logs {
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
