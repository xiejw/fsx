package cmdlog

import (
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
	logs := CmdLogs([]CmdLog{
		{
			Version:    0,
			MasterName: "test_node",
			Timestamp:  123,
			Cmd:        &Cmd{CmdNew, "", "test_file_1", 123, []byte("abc")},
		},
	})

	err := logs.Check()
	assertNoError(t, err)
}

func TestCheckFailedWithCmdLogVersion(t *testing.T) {
	logs := CmdLogs([]CmdLog{
		{
			Version:    1,
			MasterName: "test_node",
			Timestamp:  123,
			Cmd:        &Cmd{CmdNew, "", "test_file_1", 123, []byte("abc")},
		},
	})

	err := logs.Check()
	assertErrorMsg(t, err, "CmdLog version")
}

func TestCheckFailedWithMissingTimestamp(t *testing.T) {
	logs := CmdLogs([]CmdLog{
		{
			Version:    0,
			MasterName: "test_node",
			Cmd:        &Cmd{CmdNew, "", "test_file_1", 123, []byte("abc")},
		},
	})

	err := logs.Check()
	assertErrorMsg(t, err, "timestamp must be set")
}

func TestCheckFailedWithMissingMasterNode(t *testing.T) {
	logs := CmdLogs([]CmdLog{
		{
			Version:   0,
			Timestamp: 123,
			Cmd:       &Cmd{CmdNew, "", "test_file_1", 123, []byte("abc")},
		},
	})

	err := logs.Check()
	assertErrorMsg(t, err, "mastername must be set")
}
