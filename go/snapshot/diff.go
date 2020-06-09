package snapshot

type DiffResult struct {
	OnlyInLhs   bool
	OnlyInRhs   bool
	DiffItem    bool
	Item        *FileItem
	AnotherItem *FileItem // If DiffItem is true, this points to rhs
}

// Diffs the snapshots
func Diff(lhs, rhs Snapshot) []*DiffResult {
	var diff []*DiffResult

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
			} else if !isItemEqual(item, rhsItem) {
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

	return diff
}

func isItemEqual(lhsItem, rhsItem *FileItem) bool {
	if lhsItem.Size != rhsItem.Size {
		return false
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
