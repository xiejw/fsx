package snapshot

import (
	"testing"
)

func TestDiffEquality(t *testing.T) {
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

func TestDiffMoreinLhs(t *testing.T) {
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

func TestDiffMoreinRhs(t *testing.T) {
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
