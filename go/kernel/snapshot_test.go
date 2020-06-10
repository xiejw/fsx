package kernel

import (
	"reflect"
	"testing"

	"github.com/xiejw/fsx/go/cmdlog"
	"github.com/xiejw/fsx/go/snapshot"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertSnapshot(t *testing.T, expected []*snapshot.FileItem, sp snapshot.Snapshot) {
	t.Helper()
	expectedMap := make(map[string]*snapshot.FileItem)
	for _, p := range expected {
		expectedMap[p.FullPath] = p
	}

	iter := sp.Iterator()
	for item := iter.Next(); item != nil; item = iter.Next() {
		expectedItem, existed := expectedMap[item.FullPath]
		if !existed {
			t.Fatalf("expected %v in snapshot", item.FullPath)
		}

		if !reflect.DeepEqual(expectedItem, item) {
			t.Fatalf("item content mismatches: expected: %v, got: %v", expectedItem, item)

		}
		delete(expectedMap, item.FullPath)
	}

	if len(expectedMap) != 0 {
		t.Errorf("found more items than the snapshot: %v", expectedMap)
	}
}

func TestFetchSnapshotFromCmdLogs(t *testing.T) {
	cmdLogs := []cmdlog.CmdLog{
		{
			Version: 0,
			Cmd: &cmdlog.Cmd{
				cmdlog.CmdNew, "", "test_file_1", 123, []byte("abc")},
		},
	}
	sp, err := FetchSnapshotFromCmdLogs(cmdLogs)
	assertNoError(t, err)
	assertSnapshot(t, []*snapshot.FileItem{
		{"test_file_1", 123, []byte("abc")},
	}, sp)
}

func TestFetchSnapshotFromCmdLogsWithDel(t *testing.T) {
	cmdLogs := []cmdlog.CmdLog{
		{
			Version: 0,
			Cmd: &cmdlog.Cmd{
				cmdlog.CmdNew, "", "test_file_1", 123, []byte("abc")},
		},
		{
			Version: 1,
			Cmd: &cmdlog.Cmd{
				cmdlog.CmdNew, "a", "test_file_2", 456, []byte("def")},
		},
		{
			Version: 2,
			Cmd:     &cmdlog.Cmd{Type: cmdlog.CmdDel, Dir: "", FileName: "test_file_1"},
		},
	}
	sp, err := FetchSnapshotFromCmdLogs(cmdLogs)
	assertNoError(t, err)
	assertSnapshot(t, []*snapshot.FileItem{
		{"a/test_file_2", 456, []byte("def")},
	}, sp)
}
