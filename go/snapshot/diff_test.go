package snapshot

import (
	"testing"
)

// Asserts that two unordered list are same.
func assertDiffResults(t *testing.T,
	expectedInLhs, expectedInRhs, expectedDiffItem map[string]bool,
	got []*DiffResult) {

	t.Helper()

	for _, result := range got {
		if result.OnlyInLhs {
			_, existed := expectedInLhs[result.Item.FullPath]
			if !existed {
				t.Errorf("expected %v in lhs", result.Item.FullPath)
			}
			delete(expectedInLhs, result.Item.FullPath)
		} else if result.OnlyInRhs {
			_, existed := expectedInRhs[result.Item.FullPath]
			if !existed {
				t.Errorf("expected %v in rhs", result.Item.FullPath)
			}
			delete(expectedInRhs, result.Item.FullPath)
		} else if result.DiffItem {
			_, existed := expectedDiffItem[result.Item.FullPath]
			if !existed {
				t.Errorf("expected %v in diff", result.Item.FullPath)
			}
			delete(expectedDiffItem, result.Item.FullPath)
		} else {
			t.Errorf("unexpected DiffResult: %v", result)
		}
	}

	if len(expectedInLhs) != 0 {
		t.Errorf("found more items in expectedInLhs: %v", expectedInLhs)
	}
	if len(expectedInRhs) != 0 {
		t.Errorf("found more items in expectedInRhs: %v", expectedInRhs)
	}
	if len(expectedDiffItem) != 0 {
		t.Errorf("found more items in expectedDiffItem: %v", expectedDiffItem)
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

	results := Diff(lhs, rhs)
	assertDiffResults(t,
		map[string]bool{"dir/b": true, "dir/c": true},
		map[string]bool{"dir/d": true, "dir/e": true},
		map[string]bool{},
		results)
}

func TestDiffWithMetaData(t *testing.T) {
	var err error
	lhs := New()

	err = lhs.Add(&FileItem{FullPath: "dir/a", Size: 123})
	assertNoErr(t, err)
	err = lhs.Add(&FileItem{FullPath: "dir/b", Checksum: []byte{'1', '2'}})
	assertNoErr(t, err)
	err = lhs.Add(&FileItem{FullPath: "dir/f", Size: 123})
	assertNoErr(t, err)

	rhs := New()
	err = rhs.Add(&FileItem{FullPath: "dir/a", Size: 123})
	assertNoErr(t, err)
	err = rhs.Add(&FileItem{FullPath: "dir/b", Checksum: []byte{'7', '2'}})
	assertNoErr(t, err)
	err = rhs.Add(&FileItem{FullPath: "dir/f", Size: 456})
	assertNoErr(t, err)

	results := Diff(lhs, rhs)
	assertDiffResults(t,
		map[string]bool{},
		map[string]bool{},
		map[string]bool{"dir/b": true, "dir/f": true},
		results)
}
