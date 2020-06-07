package scanner

import (
	"testing"
)

func TestPrintPathAndSizeFormatter(t *testing.T) {
	formatter := NewFormatter(FormatterOption{
		PrintChecksum: false,
	})
	metadata := FileMetadata{
		BaseDir: "/abc",
		Path:    "d/e/f",
		Info: &newFileInfo{
			name:    "/f",
			is_file: true,
			size:    123,
		},
	}

	oldPrintf := printf
	calledCount := 0
	printf = func(fmt string, a ...interface{}) (n int, err error) {
		if a[0].(int64) != 123 {
			t.Errorf("File size is wrong. got %v", a[0])
		}
		if a[1].(string) != "d/e/f" {
			t.Errorf("File path is wrong.")
		}
		calledCount += 1
		return 0, nil
	}
	defer func() {
		printf = oldPrintf
	}()
	formatter(metadata)
	if calledCount != 1 {
		t.Errorf("Should be called once.")
	}
}

func TestPrintPathAndSizeFormatterSkipDir(t *testing.T) {
	formatter := NewFormatter(FormatterOption{
		PrintChecksum: false,
	})
	metadata := FileMetadata{
		BaseDir: "/abc",
		Path:    "d/e/f",
		Info: &newFileInfo{
			name:    "/f",
			is_file: false,
		},
	}

	oldPrintf := printf
	printf = func(fmt string, a ...interface{}) (n int, err error) {
		t.Fatalf("Should skip directory.")
		return 0, nil
	}
	defer func() {
		printf = oldPrintf
	}()
	formatter(metadata)
}
