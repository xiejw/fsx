package kernel

import (
	"fmt"

	"github.com/golang/glog"
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

func FetchSnapshotFromCmdLogs(cmdLogs bool) {
}
