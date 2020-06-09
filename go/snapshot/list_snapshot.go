package snapshot

// One implementation for Snaptshot based on list.
//
// Internally, a map is created to speed up the LookUp.
type ListSnapshot struct {
	array []*FileItem    // nil is the hole.
	index map[string]int // Index the path to array index.
	holes []int          // Stores the hole index in array.
}

func New() SnapshotBuilder {
	return &ListSnapshot{
		array: make([]*FileItem, 0),
		index: make(map[string]int),
		holes: make([]int, 0),
	}
}

func (sp *ListSnapshot) LookUp(fullPath string) *FileItem {
	index, exist := sp.index[fullPath]
	if !exist {
		return nil
	}
	return sp.array[index]
}

func (sp *ListSnapshot) Add(item *FileItem) error {
	fullPath := item.FullPath
	_, exist := sp.index[fullPath]
	if exist {
		return ErrItemAlreadyExist
	}

	// First find the avaiable hole.
	count := len(sp.holes)
	if count != 0 {
		holeIndex := sp.holes[count-1]
		sp.holes = sp.holes[:count-1] // Shrink it.
		sp.index[fullPath] = holeIndex
		if sp.array[holeIndex] != nil {
			panic("internal error: hole should be nil.")
		}
		sp.array[holeIndex] = item
		return nil
	}

	holeIndex := len(sp.array)
	sp.array = append(sp.array, item)
	sp.index[fullPath] = holeIndex
	return nil
}

func (sp *ListSnapshot) Delete(fullPath string) error {
	deletedIndex, exist := sp.index[fullPath]
	if !exist {
		return ErrItemNotExist
	}
	delete(sp.index, fullPath)
	sp.array[deletedIndex] = nil

	// Put the index to holes
	sp.holes = append(sp.holes, deletedIndex)
	return nil
}

type ListSnapshotIterator struct {
	array      []*FileItem
	index      int
	finalIndex int
}

func (sp *ListSnapshot) Iterator() Iterator {
	return &ListSnapshotIterator{
		array:      sp.array,
		index:      0,
		finalIndex: len(sp.array),
	}
}

func (spi *ListSnapshotIterator) Next() *FileItem {
	for {
		if spi.index >= spi.finalIndex {
			return nil
		}

		item := spi.array[spi.index]
		spi.index++
		if item == nil {
			continue
		}

		return item
	}
}
