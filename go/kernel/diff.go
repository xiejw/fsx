package kernel

import (
	"github.com/xiejw/fsx/go/cmdlog"
	"github.com/xiejw/fsx/go/snapshot"
)

// Diffs the snapshots and generates the CmdLogs.
//
// Only the Cmd in each CmdLog is filled.
func ComputeCmdLogs(newSp, oldSp snapshot.Snapshot) (cmdlog.CmdLogs, error) {
	return nil, nil
}
