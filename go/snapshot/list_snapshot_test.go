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

// Asserts that items in sp are matching the list of `expected`.
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

func TestAddDup(t *testing.T) {
	var err error
	sp := New()
	err = sp.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)

	err = sp.Add(&FileItem{FullPath: "dir/a"})
	if err != ErrItemAlreadyExist {
		t.Errorf("error mismatches.")
	}
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

func TestDeleteNotExist(t *testing.T) {
	var err error
	sp := New()

	err = sp.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)

	err = sp.Delete("dir/a")
	assertNoErr(t, err)

	err = sp.Delete("dir/a")
	if err != ErrItemNotExist {
		t.Errorf("error mismatches.")
	}
}

func TestLookUp(t *testing.T) {
	var err error
	sp := New()
	err = sp.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)

	if sp.LookUp("dir/b") != nil {
		t.Errorf("unexpected item for `dir/b`.")
	}

	if sp.LookUp("dir/a") == nil {
		t.Errorf("unexpected missing item for `dir/a`.")
	}
}

func TestComplicatedCase(t *testing.T) {
	var err error
	sp := New()

	err = sp.Add(&FileItem{FullPath: "dir/a"})
	assertNoErr(t, err)
	err = sp.Add(&FileItem{FullPath: "dir/b"})
	assertNoErr(t, err)
	err = sp.Add(&FileItem{FullPath: "dir/c"})
	assertNoErr(t, err)

	err = sp.Delete("dir/b")
	assertNoErr(t, err)
	err = sp.Delete("dir/c")
	assertNoErr(t, err)

	err = sp.Add(&FileItem{FullPath: "dir/d"})
	assertNoErr(t, err)
	err = sp.Add(&FileItem{FullPath: "dir/e"})
	assertNoErr(t, err)
	err = sp.Add(&FileItem{FullPath: "dir/f"})
	assertNoErr(t, err)

	err = sp.Delete("dir/f")
	assertNoErr(t, err)

	assertSnapshot(t, []string{"dir/a", "dir/d", "dir/e"}, sp)
}
