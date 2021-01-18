package snapshot

// Myer's diff algorithm.
func Diff(lhs, rhs *Region) int {
	return shortestEditSequence(lhs, rhs)
}

// Returns the shortest edit sequence that converts lhs into rhs.
func shortestEditSequence(lhs, rhs *Region) int {
	lItems := lhs.Items
	rItems := rhs.Items
	n := len(lItems)
	m := len(rItems)

	maxD := n + m
	offset := maxD

	var x int
	var y int

	V := make([]int, 2*maxD+1) // zeros all. in paper only V[1+offset] is needed.

	for d := 0; d <= maxD; d++ {
		for k := -d; k <= d; k += 2 {
			if k == -d || (k != d && V[k-1+offset] < V[k+1+offset]) {
				x = V[k+1+offset] // down
			} else {
				x = V[k-1+offset] + 1 // right
			}
			y = x - k

			// follow snake. different from paper, items are zero-index based.
			for x < n && y < m && lItems[x].RelPath == rItems[y].RelPath { // Use eq
				x++
				y++
			}

			V[k+offset] = x

			// check for solution
			if x >= n && y >= m {
				return d
			}
		}
	}
	return -1
}
