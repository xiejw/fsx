package cmdlogs

import (
	"bytes"
	"testing"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck(t *testing.T) {
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
