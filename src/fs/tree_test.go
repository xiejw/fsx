package fs

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/xiejw/fsx/src/clogs"
	"github.com/xiejw/fsx/src/fs/internal/scanner"
)

func TestFileItemStructSameAsCLogsFileItem(t *testing.T) {
	fi := &FileItem{Path: "afoo", Size: 123, Checksum: "0xabc"}
	clgs := &clogs.FileItem{}

	val_1 := reflect.ValueOf(fi).Elem()
	val_2 := reflect.ValueOf(clgs).Elem()

	if val_1.NumField() != val_2.NumField() {
		t.Fatalf("num fields mismatch.")
	}

	for i := 0; i < val_1.NumField(); i++ {
		if val_1.Type().Field(i).Name != val_2.Type().Field(i).Name {
			t.Errorf("field %v field name mismatch.", i)
		}
		if val_1.Field(i).Kind() != val_2.Field(i).Kind() {
			t.Errorf("field %v field type mismatch.", i)
		}
	}
}

func TestFromWalk(t *testing.T) {
	walkFn := func(baseDir string, filters []scanner.Filter, formatter scanner.Formatter) error {
		metadata := scanner.FileMetadata{
			Path: "a",
			Info: &newFileInfo{is_file: false},
		}
		formatter(metadata)

		metadata.Path = "b"
		metadata.Info = &newFileInfo{is_file: true, size: 123}
		formatter(metadata)

		metadata.Path = "c"
		metadata.Info = &newFileInfo{is_file: true, size: 456}
		formatter(metadata)

		return nil
	}

	ft, err := fromLocalFS("abc", walkFn)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if ft.BaseDir != "abc" {
		t.Errorf("unexpected basedir: %v", err)
	}

	if ft.HasChecksum {
		t.Errorf("unexpected checksum state.")
	}

	expected := []*FileItem{
		{"b", 123, ""},
		{"c", 456, ""},
	}

	if !reflect.DeepEqual(ft.Items, expected) {
		t.Errorf("unexpected items state. expected:\n%+v\ngot:\n%+v\n", expected, ft.Items)
	}
}

func TestFromCLogs(t *testing.T) {
	clgs := &clogs.CmdLogs{
		Cmds: []*clogs.CmdLog{
			{clogs.CmdNew, clogs.FileItem{"c", 5, "0xb"}, 0},
			{clogs.CmdNew, clogs.FileItem{"bc", 3, "0xc"}, 0},
			{clogs.CmdNew, clogs.FileItem{"b/c", 2, "0xb"}, 0},
			{clogs.CmdNew, clogs.FileItem{"a/b", 1, "0x1"}, 0},
			{clogs.CmdDel, clogs.FileItem{"a/b", 1, "0x1"}, 0},
			{clogs.CmdNew, clogs.FileItem{"a/b", 4, "0x1"}, 0},
			{clogs.CmdDel, clogs.FileItem{"c", 5, "0xb"}, 0},
		},
	}

	ft, err := FromCmdLogs("abc", clgs)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if ft.BaseDir != "abc" {
		t.Errorf("unexpected basedir: %v", err)
	}

	if !ft.HasChecksum {
		t.Errorf("unexpected checksum state.")
	}

	expected := []*FileItem{
		{"a/b", 4, "0x1"},
		{"b/c", 2, "0xb"},
		{"bc", 3, "0xc"},
	}

	if !reflect.DeepEqual(ft.Items, expected) {
		t.Errorf("unexpected items state. expected:\n%+v\ngot:\n%+v\n", expected, ft.Items)
	}
}

func TestDiffWoCmpChecksum(t *testing.T) {
	tree1 := &FileTree{
		HasChecksum: true,
		Items: []*FileItem{
			{"a/c", 9, "0x1"},
			{"a/b", 4, "0x1"},
			{"b/c", 2, "0xb"},
			{"bc", 3, "0xc"},
		},
	}
	tree2 := &FileTree{
		HasChecksum: false,
		Items: []*FileItem{
			{"a/b", 4, ""},
			{"b/c", 3, ""},
			{"bc", 3, ""},
			{"bcd", 6, ""},
		},
	}

	del, add, err := tree1.ConvertTo(tree2)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	expected_del := []*FileItem{
		{"a/c", 9, "0x1"},
		{"b/c", 2, "0xb"},
	}
	assertFileItemListEqual(t, "del mismatch.", expected_del, del)

	expected_add := []*FileItem{
		{"b/c", 3, ""},
		{"bcd", 6, ""},
	}
	assertFileItemListEqual(t, "add mismatch.", expected_add, add)
}

func TestDiffWCmpChecksum(t *testing.T) {
	tree1 := &FileTree{
		HasChecksum: true,
		Items: []*FileItem{
			{"a/b", 4, "0x1"},
			{"b/c", 2, "0xb"},
		},
	}
	tree2 := &FileTree{
		HasChecksum: true,
		Items: []*FileItem{
			{"a/b", 4, "0x2"},
			{"b/c", 2, "0xb"},
		},
	}

	del, add, err := tree1.ConvertTo(tree2)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	expected_del := []*FileItem{
		{"a/b", 4, "0x1"},
	}
	assertFileItemListEqual(t, "del mismatch.", expected_del, del)

	expected_add := []*FileItem{
		{"a/b", 4, "0x2"},
	}
	assertFileItemListEqual(t, "add mismatch.", expected_add, add)
}

// -------------------------------------------------------------------------------------------------
// helper
// -------------------------------------------------------------------------------------------------

type newFileInfo struct {
	name    string
	is_file bool
	mode    int
	size    int64
}

func (f *newFileInfo) Name() string {
	return f.name
}

func (f *newFileInfo) Size() int64 {
	return f.size
}
func (f *newFileInfo) Mode() os.FileMode {
	return os.FileMode(f.mode)
}
func (f *newFileInfo) ModTime() time.Time {
	return time.Now()
}
func (f *newFileInfo) IsDir() bool {
	return !f.is_file
}
func (f *newFileInfo) Sys() interface{} {
	return nil
}

func assertFileItemListEqual(t *testing.T, err_msg string, expected, got []*FileItem) {
	if !reflect.DeepEqual(expected, got) {
		err_msg += "\nexpected:\n"
		for _, fi := range expected {
			err_msg += fmt.Sprintf("  %+v\n", fi)
		}
		err_msg += "got:\n"
		for _, fi := range got {
			err_msg += fmt.Sprintf("  %+v\n", fi)
		}
		t.Fatalf(err_msg)
	}
}
