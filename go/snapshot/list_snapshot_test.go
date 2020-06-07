package snapshot

import (
	"testing"
)

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unepxected error: %v", err)
	}
}

func assertSnapshot(t *testing.T, expected []string, sp Snapshot) {
	t.Helper()
	expectedMap := make(map[string]bool)
	for _, p := range expected {
		expectedMap[p] = true
	}

	iter := sp.Iterator()
	for {
		item := iter.Next()
		if item == nil {
			break
		}

		_, existed := expectedMap[item.FullPath]
		if !existed {
			t.Errorf("expected %v in snapshot", item.FullPath)
		}
		delete(expectedMap, item.FullPath)
	}

	if len(expectedMap) != 0 {
		t.Errorf("found more items than the snapshot: %v", expectedMap)
	}
}

func TestAdd(t *testing.T) {
	var err error
	sp := New()
	err = sp.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)

	assertSnapshot(t, []string{"dir/a"}, sp)
}

func TestDelete(t *testing.T) {
	var err error
	sp := New()

	err = sp.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)

	err = sp.Delete("dir/a")
	assertNoErr(t, err)

	assertSnapshot(t, []string{}, sp)
}
