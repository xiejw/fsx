package snapshot

func Diff(lhs, rhs Snapshot) ([]*FileItem, []*FileItem) {
	var onlyInLhs, onlyInRhs []*FileItem

	{
		iter := lhs.Iterator()
		for {
			item := iter.Next()
			if item == nil {
				break
			}

			if rhs.LookUp(item.FullPath) == nil {
				onlyInLhs = append(onlyInLhs, item)
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
				onlyInRhs = append(onlyInRhs, item)
			}
		}
	}

	return onlyInLhs, onlyInRhs
}
