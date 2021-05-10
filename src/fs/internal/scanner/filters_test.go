package scanner

import (
	"os"
	"testing"
	"time"
)

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

var _fakeFileName = "/tmp"

func TestCommonFilter(t *testing.T) {
	commonFilter := NewCommonFilter(nil)

	regularFile := &newFileInfo{
		name:    "/tmp",
		is_file: true,
		mode:    0,
	}
	if commonFilter("/tmp", regularFile) {
		t.Fatal("Regualr file should never be skipped.")
	}

	blacklistFileInfos := make([]*newFileInfo, 0)

	blacklistFileInfos = append(blacklistFileInfos, &newFileInfo{
		name:    ".git",
		is_file: true,
	})

	blacklistFileInfos = append(blacklistFileInfos, &newFileInfo{
		name:    ".git",
		is_file: false,
	})

	blacklistFileInfos = append(blacklistFileInfos, &newFileInfo{
		name:    ".DS_Store",
		is_file: true,
	})

	blacklistFileInfos = append(blacklistFileInfos, &newFileInfo{
		name:    "._Image",
		is_file: true,
	})

	blacklistFileInfos = append(blacklistFileInfos, &newFileInfo{
		name:    "go.swp",
		is_file: true,
	})

	for _, fileInfo := range blacklistFileInfos {
		if !commonFilter(fileInfo.name, fileInfo) {
			t.Fatalf("%v (is_file: %v) should be skipped.", fileInfo.name, fileInfo.is_file)
		}
	}
}

func TestCommonFilterWithExtsToSkip(t *testing.T) {
	commonFilter := NewCommonFilter(nil)
	regularFile := &newFileInfo{
		name:    "tmp.ini",
		is_file: true,
	}
	if commonFilter("tmp.ini", regularFile) {
		t.Fatal("Regualr file should never be skipped.")
	}

	commonFilter = NewCommonFilter([]string{"ini"})
	if !commonFilter("tmp.ini", regularFile) {
		t.Fatal("File with .ini extension should be skipped.")
	}
}

func TestCounterFilter(t *testing.T) {
	counterFilter := NewCounterFilter(2, false /*fileOnly*/)

	f1 := &newFileInfo{
		name:    _fakeFileName,
		is_file: true,
	}
	f2 := &newFileInfo{
		name:    _fakeFileName,
		is_file: false,
	}
	f3 := &newFileInfo{
		name:    _fakeFileName,
		is_file: true,
	}
	if counterFilter(_fakeFileName, f1) {
		t.Fatal("f1 should never be skipped.")
	}
	if counterFilter(_fakeFileName, f2) {
		t.Fatal("f2 should never be skipped.")
	}
	if !counterFilter(_fakeFileName, f3) {
		t.Fatal("f3 should be skipped.")
	}
}

func TestCounterFilterWithFileOnly(t *testing.T) {
	counterFilter := NewCounterFilter(1, true /*fileOnly*/)

	f1 := &newFileInfo{
		name:    _fakeFileName,
		is_file: false,
	}
	f2 := &newFileInfo{
		name:    _fakeFileName,
		is_file: true,
	}
	f3 := &newFileInfo{
		name:    _fakeFileName,
		is_file: true,
	}
	if counterFilter(_fakeFileName, f1) {
		t.Fatal("f1 should never be skipped.")
	}
	if counterFilter(_fakeFileName, f2) {
		t.Fatal("f2 should never be skipped.")
	}
	if !counterFilter(_fakeFileName, f3) {
		t.Fatal("f3 should be skipped.")
	}
}
