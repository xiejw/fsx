package snapshot

import (
	"testing"
)

// Asserts that two unordered list are same.
func assertFileItemList(t *testing.T, expected []string, got []*FileItem) {
	t.Helper()
	expectedMap := make(map[string]bool)
	for _, p := range expected {
		expectedMap[p] = true
	}

	for _, item := range got {
		_, existed := expectedMap[item.FullPath]
		if !existed {
			t.Errorf("expected %v in file item list", item.FullPath)
		}
		delete(expectedMap, item.FullPath)
	}

	if len(expectedMap) != 0 {
		t.Errorf("found more items than the file item list: %v", expectedMap)
	}
}

func TestDiff(t *testing.T) {
	var err error
	lhs := New()

	err = lhs.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)
	err = lhs.Add(&FileItem{FullPath: "dir/f"})
	assertNoErr(t, err)
	err = lhs.Add(&FileItem{FullPath: "dir/b"})
	assertNoErr(t, err)
	err = lhs.Add(&FileItem{FullPath: "dir/c"})
	assertNoErr(t, err)

	rhs := New()
	err = rhs.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)
	err = rhs.Add(&FileItem{FullPath: "dir/d"})
	assertNoErr(t, err)
	err = rhs.Add(&FileItem{FullPath: "dir/e"})
	assertNoErr(t, err)
	err = rhs.Add(&FileItem{FullPath: "dir/f"})
	assertNoErr(t, err)

	lhsfis, rhsfis := Diff(lhs, rhs)
	assertFileItemList(t, []string{"dir/b", "dir/c"}, lhsfis)
	assertFileItemList(t, []string{"dir/d", "dir/e"}, rhsfis)
}
