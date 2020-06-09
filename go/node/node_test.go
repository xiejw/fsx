package node

import (
	_ "bytes"
	_ "reflect"
	"strings"
	"testing"

	"github.com/xiejw/fsx/go/cmdlog"
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
	state := &Node{
		Name:        "test_node",
		NextVersion: 1,
		CmdLogs: []cmdlog.CmdLog{
			{
				Version:    0,
				MasterName: "test_node",
				Cmd:        &cmdlog.Cmd{cmdlog.CmdNew, "", "test_file_1", 123, []byte("abc")},
			},
		},
	}

	err := state.Check()
	assertNoError(t, err)
}

func TestCheckFailedWithNextVersion(t *testing.T) {
	state := &Node{
		Name:        "test_node",
		NextVersion: 0,
		CmdLogs: []cmdlog.CmdLog{
			{
				Version:    0,
				MasterName: "test_node",
				Cmd:        &cmdlog.Cmd{cmdlog.CmdNew, "", "test_file_1", 123, []byte("abc")},
			},
		},
	}

	err := state.Check()
	assertErrorMsg(t, err, "NextVersion")
}

func TestCheckFailedWithCmdLogVersion(t *testing.T) {
	state := &Node{
		Name:        "test_node",
		NextVersion: 1,
		CmdLogs: []cmdlog.CmdLog{
			{
				Version:    1,
				MasterName: "test_node",
				Cmd:        &cmdlog.Cmd{cmdlog.CmdNew, "", "test_file_1", 123, []byte("abc")},
			},
		},
	}

	err := state.Check()
	assertErrorMsg(t, err, "CmdLog version")
}

// func TestMarshal(t *testing.T) {
// 	state := &Node{
// 		Name:        "test_node",
// 		NextVersion: 1,
// 		CmdLogs: []CmdLog{
// 			{
// 				Version:    0,
// 				MasterName: "test_node",
// 				Cmd:        &Cmd{CmdNew, "", "test_file_1", 123, []byte("abc")},
// 			},
// 		},
// 	}
//
// 	var buf bytes.Buffer
// 	err := state.Marshal(&buf)
// 	assertNoError(t, err)
//
// 	got := buf.String()
// 	expected := `{
//   "Name": "test_node",
//   "NextVersion": 1,
//   "IsMaster": false,
//   "CmdLogs": [
//     {
//       "Version": 0,
//       "MasterName": "test_node",
//       "Cmd": {
//         "Type": 1,
//         "Dir": "",
//         "FileName": "test_file_1",
//         "Size": 123,
//         "Checksum": "YWJj"
//       }
//     }
//   ]
// }`
// 	if got != expected {
// 		t.Fatalf("marshal content mismatches. expected: %v, got: %v",
// 			expected, got)
// 	}
// }
//
// func TestUnmarshal(t *testing.T) {
// 	state := &Node{
// 		Name:        "test_node",
// 		NextVersion: 1,
// 		CmdLogs: []CmdLog{
// 			{
// 				Version:    0,
// 				MasterName: "test_node",
// 				Cmd:        &Cmd{CmdNew, "", "test_file_1", 123, []byte("abc")},
// 			},
// 		},
// 	}
//
// 	var buf bytes.Buffer
// 	err := state.Marshal(&buf)
// 	assertNoError(t, err)
//
// 	got, err := Unmarshal(buf.Bytes())
// 	assertNoError(t, err)
//
// 	if !reflect.DeepEqual(got, state) {
// 		t.Fatalf("unmarshal content mismatches. expected: %v, got: %v",
// 			state, got)
// 	}
// }
