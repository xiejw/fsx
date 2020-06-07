package cmdlogs

import (
	"bytes"
	"strings"
	"testing"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertErrorMsg(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error")
	}

	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("expected to get error message: %v, got: %v",
			msg, err)
	}
}

func TestCheckPassed(t *testing.T) {
	state := &NodeState{
		Name:        "test_node",
		NextVersion: 1,
		CmdLogs: []CmdLog{
			{
				Version:    0,
				MasterName: "test_node",
				Cmd:        &Cmd{CmdNew, "", "test_file_1", []byte("abc")},
			},
		},
	}

	err := state.Check()
	assertNoError(t, err)
}

func TestCheckFailedWithNextVersion(t *testing.T) {
	state := &NodeState{
		Name:        "test_node",
		NextVersion: 0,
		CmdLogs: []CmdLog{
			{
				Version:    0,
				MasterName: "test_node",
				Cmd:        &Cmd{CmdNew, "", "test_file_1", []byte("abc")},
			},
		},
	}

	err := state.Check()
	assertErrorMsg(t, err, "NextVersion")
}

func TestMarshal(t *testing.T) {
	state := &NodeState{
		Name:        "test_node",
		NextVersion: 1,
		CmdLogs: []CmdLog{
			{
				Version:    0,
				MasterName: "test_node",
				Cmd:        &Cmd{CmdNew, "", "test_file_1", []byte("abc")},
			},
		},
	}

	var buf bytes.Buffer
	err := state.Marshal(&buf)
	assertNoError(t, err)
}
