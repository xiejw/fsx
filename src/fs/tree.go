package fs

import (
	"sort"

	"github.com/xiejw/fsx/src/clogs"
	"github.com/xiejw/fsx/src/errors"
	"github.com/xiejw/fsx/src/fs/internal/scanner"
)

// -------------------------------------------------------------------------------------------------
// public
// -------------------------------------------------------------------------------------------------

// same as clogs.FileItem
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

// creates a FileTree by walking the baseDir.
func FromLocalFS(baseDir string) (*FileTree, error) {
	return fromLocalFS(baseDir, scanner.Walk)
}

// creates a FileTree by replaying cmds in CmdLogs.
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

// steps to convert from src to dst by deleting items in del first, followed by adding items in add.
func (src *FileTree) ConvertTo(dst *FileTree) (del []*FileItem, add []*FileItem, err error) {
	cmp_checksum := false
	if src.HasChecksum && dst.HasChecksum {
		cmp_checksum = true
	}
	diffs, err := diff(src.Items, dst.Items, cmp_checksum)
	if err != nil {
		return
	}

	for _, r := range diffs {
		// for all cases defined by diffChange, the rule is same.
		if r.lhs != nil {
			del = append(del, r.lhs)
		}
		if r.rhs != nil {
			add = append(add, r.rhs)
		}
	}
	return
}

// -------------------------------------------------------------------------------------------------
// impl
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

type diffChange struct {
	lhs *FileItem // existed only in lhs. if both set, they are diff.
	rhs *FileItem // existed only in rhs.
}

func diff(lhs, rhs []*FileItem, checksum bool) ([]*diffChange, error) {
	return nil, nil
}
