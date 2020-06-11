package snapshot

// One implementation for Snaptshot based on list.
//
// Internally, a map is created to speed up the LookUp.
type ListSnapshot struct {
	size        uint64         // size of items.
	array       []*FileItem    // nil is the hole.
	index       map[string]int // Index the path to array index.
	holes       []int          // Stores the hole index in array.
	hasChecksum bool           // The state of hasChecksum for all items.
}

func New() SnapshotBuilder {
	return &ListSnapshot{
		array: make([]*FileItem, 0),
		index: make(map[string]int),
		holes: make([]int, 0),
	}
}

func (sp *ListSnapshot) HasChecksum() bool {
	return sp.hasChecksum
}

func (sp *ListSnapshot) Size() uint64 {
	return sp.size
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

	// Check the invariance of the `hasChecksum`.
	itemHasChecksum := len(item.Checksum) != 0
	if sp.size == 0 {
		sp.hasChecksum = itemHasChecksum
	} else if sp.hasChecksum != itemHasChecksum {
		return ErrChecksumMismatch
	}

	// Find the avaiable hole first.
	sp.size++
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
	sp.size--
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
