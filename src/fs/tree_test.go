package fs

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/xiejw/fsx/src/clogs"
	"github.com/xiejw/fsx/src/fs/scanner"
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

func TestWalk(t *testing.T) {
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
