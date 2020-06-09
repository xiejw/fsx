package kernel

import (
	"fmt"
	"path"

	"github.com/golang/glog"
	"github.com/xiejw/fsx/go/cmdlog"
	"github.com/xiejw/fsx/go/fs"
	"github.com/xiejw/fsx/go/snapshot"
)

func FetchSnapshotFromFileTree(baseDir string) (snapshot.Snapshot, error) {
	sp := snapshot.New()

	var errDuringWalk error

	err := fs.Walk(".", func(walkResult *fs.WalkResult) {
		if walkResult.Err != nil {
			errDuringWalk = walkResult.Err
			glog.Errorf("unexpected error: %v.", errDuringWalk)
			return
		}
		sp.Add(&snapshot.FileItem{
			FullPath: walkResult.Path,
			Size:     walkResult.Size,
		})
	})

	if err != nil {
		return nil, fmt.Errorf("failed to start walking the file tree: %w", err)
	}

	if errDuringWalk != nil {
		return nil, fmt.Errorf("failed to read file item under the file tree: %w", err)
	}

	return sp, nil
}

func FetchSnapshotFromCmdLogs(cmdLogs []cmdlog.CmdLog) (snapshot.Snapshot, error) {
	var err error
	sp := snapshot.New()

	for i, log := range cmdLogs {
		cmd := log.Cmd
		switch cmd.Type {
		case cmdlog.CmdNew:
			item := &snapshot.FileItem{
				FullPath: path.Join(cmd.Dir, cmd.FileName),
				Size:     cmd.Size,
				Checksum: cmd.Checksum,
			}
			err = sp.Add(item)
			if err != nil {
				return nil, err
			}
		case cmdlog.CmdDel:
			fullPath := path.Join(cmd.Dir, cmd.FileName)
			err = sp.Delete(fullPath)
			if err != nil {
				return nil, err
			}

		default:
			return nil, fmt.Errorf("unsupported cmdlog type at %v: %v", i, cmd.Type)
		}
	}

	return sp, nil
}
