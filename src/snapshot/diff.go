package snapshot

import "fmt"

// Myer's diff algorithm.
func Diff(lhs, rhs *Region) int {
	lItems := lhs.Items
	rItems := rhs.Items
	n := len(lItems)
	m := len(rItems)

	maxD := n + m
	offset := maxD

	var x int
	var y int

	V := make([]int, 2*maxD+1)

	for d := 0; d <= maxD; d++ {
		for k := -d; k <= d; k += 2 {
			// down or right?
			if k == -d || (k != d && V[k-1+offset] < V[k+1+offset]) {
				x = V[k+1+offset]
			} else {
				x = V[k-1+offset] + 1
			}
			y = x - k

			// follow snake
			for x < n && y < m && lItems[x+1].RelPath == rItems[y+1].RelPath { // Use eq
				x++
				y++
			}

			V[k+offset] = x

			// check for solution
			if x >= n && y >= m {
				fmt.Printf("len is %v\n", d)
				return d
			}
		}
	}
	return -1
}
