package fs

import (
	"fmt"
	"io"
	"path"
	"sort"

	"github.com/xiejw/fsx/src/clogs"
	"github.com/xiejw/fsx/src/errors"
	"github.com/xiejw/fsx/src/fs/internal/crypto"
	"github.com/xiejw/fsx/src/fs/internal/scanner"
)

// -------------------------------------------------------------------------------------------------
// public
// -------------------------------------------------------------------------------------------------

// FileItem should have same structure as clogs.FileItem.
type FileItem struct {
	Path     string // relative to domain.
	Size     int64  // file size.
	Checksum string // optionel checksum.
}

type FileTree struct {
	BaseDir     string      // domain path.
	HasChecksum bool        // if false, all items do not have checksum.
	Items       []*FileItem // alphabetically sorted items.
}

// FromLocalFS creates a FileTree by walking the baseDir.
func FromLocalFS(baseDir string, checksum bool) (*FileTree, error) {
	ft, err := fromLocalFS(baseDir, scanner.Walk)
	if err != nil {
		return nil, err
	}

	if checksum {
		for _, fi := range ft.Items {
			absPath := path.Join(baseDir, fi.Path)
			sum, err := crypto.Sha256Sum(absPath)
			if err != nil {
				return nil, errors.WrapNote(err, "failed to fetch checksum for file: %v", absPath)
			}
			fi.Checksum = fmt.Sprintf("%x", sum)
		}
		ft.HasChecksum = true
	}

	return ft, nil
}

// FromCmdLogs creates a FileTree by replaying cmds in CmdLogs.
func FromCmdLogs(baseDir string, clgs *clogs.CmdLogs) (*FileTree, error) {
	items, err := fromCmdLogs(baseDir, clgs)
	if err != nil {
		return nil, err
	}
	return &FileTree{
		BaseDir:     baseDir,
		HasChecksum: true,
		Items:       items,
	}, nil
}

type DiffResult struct {
	Del []*FileItem
	Add []*FileItem
}

// ConvertTo produces the steps to convert from src to dst by deleting items in del first, followed
// by adding items in add.  BaseDir is ignored during comparision. Checksum will be checked if and
// only if HasChecksum are true for both src and dst.
func (src *FileTree) ConvertTo(dst *FileTree) (*DiffResult, error) {
	// only compare checksum if both exist. useless in other caess.
	cmp_checksum := src.HasChecksum && dst.HasChecksum

	diffs, err := diff(src.Items, dst.Items, cmp_checksum)
	if err != nil {
		return nil, err
	}

	var add, del []*FileItem

	for _, r := range diffs {
		// for all cases defined by diffChange, the rule is same.
		if r.lhs != nil {
			del = append(del, r.lhs)
		}
		if r.rhs != nil {
			add = append(add, r.rhs)
		}
	}
	return &DiffResult{
		Del: del,
		Add: add,
	}, nil
}

// Print prints the DiffResult to writer.
func (diff *DiffResult) Print(w io.Writer) {
	add := diff.Add
	del := diff.Del

	fmt.Fprintf(w, "Del %v items\n", len(del))
	if len(del) == 0 {
		fmt.Fprintf(w, "  (empty)\n")
	}
	for _, it := range del {
		fmt.Fprintf(w, "  - %10d %v\n", it.Size, it.Path)
	}

	fmt.Printf("Add %v items\n", len(add))
	if len(add) == 0 {
		fmt.Fprintf(w, "  (empty)\n")
	}
	for _, it := range add {
		fmt.Fprintf(w, "  + %10d %v\n", it.Size, it.Path)
	}
}

// -------------------------------------------------------------------------------------------------
// impl (factory methods to create FileTree)
// -------------------------------------------------------------------------------------------------

func fromLocalFS(baseDir string, walkFn func(baseDir string, filters []scanner.Filter, formatter scanner.Formatter) error) (*FileTree, error) {
	ft := &FileTree{
		BaseDir:     baseDir,
		HasChecksum: false, // no checksum for perf.
		Items:       nil,
	}

	items := make([]*FileItem, 0)

	cb := func(metadata scanner.FileMetadata) {
		info := metadata.Info
		if info.IsDir() {
			return
		}
		items = append(items, &FileItem{
			Path: metadata.Path,
			Size: info.Size(),
		})
	}

	err := walkFn(baseDir, []scanner.Filter{scanner.NewCommonFilter(nil)}, cb)
	if err != nil {
		return nil, errors.WrapNote(err, "failed to create FileTree by walking local fs.")
	}

	ft.Items = items

	return ft, nil
}

// conform sort pkg
type fileItems []*FileItem

func (a fileItems) Len() int           { return len(a) }
func (a fileItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileItems) Less(i, j int) bool { return a[i].Path < a[j].Path }

func fromCmdLogs(baseDir string, clgs *clogs.CmdLogs) ([]*FileItem, error) {
	maps := make(map[string]*FileItem, len(clgs.Cmds))

	// stage 1: replay all cmds.
	for i, cmd := range clgs.Cmds {
		cfi := cmd.FileItem
		path := cfi.Path
		if path == "" {
			return nil, errors.New("the %v-th item should be empty path.", i)
		}

		item, existed := maps[path]

		switch cmd.Kind {
		case clogs.CmdNew:
			if existed {
				return nil, errors.New("the %v-th item should NOT be existed by replaying. path: %v", i, path)
			}

			fi := &FileItem{
				Path:     path,
				Size:     cfi.Size,
				Checksum: cfi.Checksum,
			}
			maps[path] = fi

		case clogs.CmdDel:
			if !existed {
				return nil, errors.New("the %v-th item should be existed by replaying. path: %v", i, path)
			}

			if item.Size != cfi.Size {
				return nil, errors.New("the %v-th item signature does not match.", i)
			}
			if item.Checksum != cfi.Checksum {
				return nil, errors.New("the %v-th item signature does not match.", i)
			}

			delete(maps, path)

		default:
			return nil, errors.New("unknown clogs Cmd Kind: %v", cmd.Kind)
		}
	}

	// stage 2: add results into items and then sort.
	items := make([]*FileItem, 0, len(maps))
	for _, fi := range maps {
		items = append(items, fi)
	}

	sort.Sort(fileItems(items))

	return items, nil
}

// -------------------------------------------------------------------------------------------------
// impl (diff FileTrees)
// -------------------------------------------------------------------------------------------------

type diffChange struct {
	lhs *FileItem // existed only in lhs. if both set, they are diff.
	rhs *FileItem // existed only in rhs. if both set, they are diff.
}

func diff(lhs, rhs []*FileItem, cmp_checksum bool) ([]*diffChange, error) {
	// A naive algorithrm. Can switch to myer's diff algorithm later, if it is too slow.

	maps_lhs := make(map[string]*FileItem, len(lhs))
	maps_rhs := make(map[string]*FileItem, len(rhs))
	df := make([]*diffChange, 0)

	// Stage 1: Push all items from rhs to maps_rhs.
	for _, fi := range rhs {
		maps_rhs[fi.Path] = fi
	}
	if len(maps_rhs) != len(rhs) {
		return nil, errors.New("FileTree has duplicated paths, which is not allowed.")
	}

	// Stage 2: Two tasks in one loop
	//   1. Push all items from lhs to maps_lhs.
	//   2. Check see whether lhs has any diff agains maps_rhs.
	for _, fi := range lhs {
		path := fi.Path
		maps_lhs[path] = fi

		fi_rhs, existed := maps_rhs[path]
		if !existed {
			df = append(df, &diffChange{lhs: fi})
			continue
		}

		// Check they are same.
		if fi.Size != fi_rhs.Size {
			df = append(df, &diffChange{lhs: fi, rhs: fi_rhs})
			continue
		}

		if cmp_checksum && fi.Checksum != fi_rhs.Checksum {
			df = append(df, &diffChange{lhs: fi, rhs: fi_rhs})
			continue
		}
	}

	if len(maps_lhs) != len(lhs) {
		return nil, errors.New("FileTree has duplicated paths, which is not allowed.")
	}

	// Stage 3: Check see whether rhs has any diff agains maps_lhs.
	for _, fi := range rhs {
		path := fi.Path
		_, existed := maps_lhs[path]
		if !existed {
			df = append(df, &diffChange{rhs: fi})
		}
		// no need to check diff as in stage 2, all are found.
	}
	return df, nil
}
