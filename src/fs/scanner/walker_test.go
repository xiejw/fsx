package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func assertStringEqual(a, b string, t *testing.T, message string, c ...interface{}) {
	if a != b {
		t.Errorf("%v. \nExpected: %v\nGot: %v", fmt.Sprintf(message, c...), a, b)
	}
}

func assertFalse(a bool, t *testing.T, message string, c ...interface{}) {
	if a {
		t.Errorf("%v. \nExpected: False\nGot: True", fmt.Sprintf(message, c...))
	}
}

func assertTrue(a bool, t *testing.T, message string, c ...interface{}) {
	if !a {
		t.Errorf("%v. \nExpected: True\nGot: False", fmt.Sprintf(message, c...))
	}
}

func setUpFilePathWalk(t *testing.T, baseDir string, callbackFn func(filepath.WalkFunc)) func() {
	oldFilePathWalk := filePathWalk
	filePathWalk = func(root string, fn filepath.WalkFunc) error {
		if root != baseDir {
			t.Errorf("Unexpected root")
		}
		callbackFn(fn)
		return nil
	}

	return func() {
		filePathWalk = oldFilePathWalk
	}
}

func TestWalkWithFormatter(t *testing.T) {
	baseDir := "/abc"

	recoverFn := setUpFilePathWalk(t, baseDir, func(walkfn filepath.WalkFunc) {
		walkfn("/abc/d/e", &newFileInfo{name: "/abc/d/e", is_file: true}, nil)
		walkfn("/abc/d/f", &newFileInfo{name: "/abc/d/f", is_file: false}, nil)
	})
	defer recoverFn()

	calledCount := 0
	formater := func(metadata FileMetadata) {
		calledCount += 1
		switch calledCount {
		case 1:
			assertStringEqual(baseDir, metadata.BaseDir, t, "Basedir should be same.")
			assertStringEqual("d/e", metadata.Path, t, "Path should be same.")
			assertFalse(metadata.Info.IsDir(), t, "Should be dir.")
		case 2:
			assertStringEqual(baseDir, metadata.BaseDir, t, "Basedir should be same.")
			assertStringEqual("d/f", metadata.Path, t, "Path should be same.")
			assertTrue(metadata.Info.IsDir(), t, "Should be dir.")
		default:
			t.Errorf("Beyond expected call count.")
		}
	}
	walkImpl(baseDir, nil, formater)
	if calledCount != 2 {
		t.Errorf("Should be called twice.")
	}
}

func TestWalkWithFilter(t *testing.T) {
	baseDir := "/abc"

	recoverFn := setUpFilePathWalk(t, baseDir, func(walkfn filepath.WalkFunc) {
		walkfn("/abc/d/e", &newFileInfo{name: "/abc/d/e"}, nil)
		walkfn("/abc/d/f", &newFileInfo{name: "/abc/d/f"}, nil)
	})
	defer recoverFn()

	mockFilter := func(path string, info os.FileInfo) bool {
		return "/abc/d/e" == path // d/e should be filtered.
	}

	calledCount := 0
	formater := func(metadata FileMetadata) {
		calledCount += 1
		assertStringEqual("d/f", metadata.Path, t, "Path should be d/f.")
	}
	walkImpl(baseDir, []Filter{mockFilter}, formater)
	if calledCount != 1 {
		t.Errorf("Should be called once.")
	}
}
