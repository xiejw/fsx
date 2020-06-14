package snapshot

type DiffResult struct {
	OnlyInLhs   bool
	OnlyInRhs   bool
	DiffItem    bool
	Item        *FileItem
	AnotherItem *FileItem // If DiffItem is true, this points to rhs
}

// Diffs the snapshots.
//
// This method compares Checksum field of Snapshot if and only if both lhs and rhs HasChecksum.  To
// ensure this is not ignored, the second return value reminds the call site this fact.
func Diff(lhs, rhs Snapshot) ([]*DiffResult, bool) {
	var diff []*DiffResult

	compareChecksum := lhs.HasChecksum() && rhs.HasChecksum()

	{
		iter := lhs.Iterator()
		for {
			item := iter.Next()
			if item == nil {
				break
			}

			// We only check diff in this pass.
			rhsItem := rhs.LookUp(item.FullPath)
			if rhsItem == nil {
				diff = append(diff, &DiffResult{
					OnlyInLhs: true,
					Item:      item,
				})
			} else if !isItemEqual(item, rhsItem, compareChecksum) {
				diff = append(diff, &DiffResult{
					DiffItem:    true,
					Item:        item,
					AnotherItem: rhsItem,
				})
			}
		}
	}

	{
		iter := rhs.Iterator()
		for {
			item := iter.Next()
			if item == nil {
				break
			}

			if lhs.LookUp(item.FullPath) == nil {
				diff = append(diff, &DiffResult{
					OnlyInRhs: true,
					Item:      item,
				})
			}
		}
	}

	return diff, compareChecksum
}

func isItemEqual(lhsItem, rhsItem *FileItem, compareChecksum bool) bool {
	if lhsItem.Size != rhsItem.Size {
		return false
	}

	if !compareChecksum {
		return true
	}

	lhs := lhsItem.Checksum
	rhs := rhsItem.Checksum

	if lhs == nil && rhs == nil {
		return true
	}
	if (lhs == nil) != (rhs == nil) {
		return false
	}

	size := len(lhs)
	if size != len(rhs) {
		return false
	}
	for i := 0; i < size; i++ {
		if lhs[i] != rhs[i] {
			return false
		}
	}
	return true
}
