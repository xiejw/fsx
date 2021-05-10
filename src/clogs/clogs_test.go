package clogs

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLogs(t *testing.T) {
	clogs := &CmdLogs{}
	if len(clogs.Cmds) != clogs.VersionID {
		t.Errorf("version ID mismatch.")
	}
}

func TestCLogsToLines(t *testing.T) {
	ep := int64(1620609633)
	clgs := []*CmdLog{
		{CmdNew, FileItem{"c", 5, "0xb"}, ep},
		{CmdNew, FileItem{"bc", 3, "0xc"}, ep + 1},
		{CmdNew, FileItem{"b/c", 2, "0xb"}, ep + 2},
		{CmdNew, FileItem{"a/b", 1, "0x1"}, ep + 2},
		{CmdDel, FileItem{"a/b", 1, "0x1"}, ep + 3},
		{CmdNew, FileItem{"a/b", 4, "0x1"}, ep + 4},
		{CmdDel, FileItem{"c", 5, "0xb"}, ep + 5},
	}

	var buf bytes.Buffer
	buf.WriteString("\n")
	for _, clg := range clgs {
		buf.WriteString(ToOneLine(clg))
		buf.WriteString("\n")
	}

	expected := `
+            5 0xb 1620609633 c
+            3 0xc 1620609634 bc
+            2 0xb 1620609635 b/c
+            1 0x1 1620609635 a/b
-            1 0x1 1620609636 a/b
+            4 0x1 1620609637 a/b
-            5 0xb 1620609638 c
`
	got := buf.String()
	if expected != got {
		t.Errorf("string mismatch.\n  expected:\n%v\n  got:\n%v\n", expected, got)
	}
}

func TestCLogsFromLines(t *testing.T) {
	r := strings.NewReader(`
+            5 0xb 1620609633 c
+            3 0xc 1620609634 bc
+            2 0xb 1620609635 b/c
+            1 0x1 1620609635 a/b
-            1 0x1 1620609636 a/b
+            4 0x1 1620609637 a/b
-            5 0xb 1620609638 c
`)
	_, err := FromLines(r)
	t.Errorf("%v", err)
}
