package fs

import (
	"reflect"
	"testing"

	"github.com/xiejw/fsx/src/clogs"
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
