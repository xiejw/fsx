package clogs

import "testing"

func TestCLogs(t *testing.T) {
	clogs := &CmdLogs{}
	if len(clogs.Cmds) != clogs.VersionID {
		t.Errorf("version ID mismatch.")
	}
}
