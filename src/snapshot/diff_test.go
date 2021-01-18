package snapshot

import (
	"testing"
)

func TestDiffLenEquality(t *testing.T) {
	lhs := Region{
		Items: []*FileItem{
			&FileItem{
				RelPath: "a",
			},
			&FileItem{
				RelPath: "b",
			},
		},
	}

	rhs := Region{
		Items: []*FileItem{
			&FileItem{
				RelPath: "a",
			},
			&FileItem{
				RelPath: "b",
			},
		},
	}

	if Diff(&lhs, &rhs) != 0 {
		t.Errorf("expect no diff.")
	}
}

func TestDiffLenMoreinLhs(t *testing.T) {
	lhs := Region{
		Items: []*FileItem{
			&FileItem{
				RelPath: "a",
			},
			&FileItem{
				RelPath: "c",
			},
			&FileItem{
				RelPath: "b",
			},
			&FileItem{
				RelPath: "d",
			},
		},
	}

	rhs := Region{
		Items: []*FileItem{
			&FileItem{
				RelPath: "a",
			},
			&FileItem{
				RelPath: "b",
			},
		},
	}

	if Diff(&lhs, &rhs) != 2 {
		t.Errorf("expect 2 diff.")
	}
}

func TestDiffLenMoreinRhs(t *testing.T) {
	lhs := Region{
		Items: []*FileItem{
			&FileItem{
				RelPath: "a",
			},
			&FileItem{
				RelPath: "b",
			},
		},
	}

	rhs := Region{
		Items: []*FileItem{
			&FileItem{
				RelPath: "a",
			},
			&FileItem{
				RelPath: "c",
			},
			&FileItem{
				RelPath: "b",
			},
			&FileItem{
				RelPath: "d",
			},
		},
	}

	if Diff(&lhs, &rhs) != 2 {
		t.Errorf("expect 2 diff.")
	}
}
