package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/xiejw/lunar"

	"github.com/xiejw/fsx/go/fs"
	"github.com/xiejw/fsx/go/snapshot"
)

func main() {
	lunar.Init(true /*parseFlag*/)
	defer lunar.FinishUp()

	sp := snapshot.New()

	err := fs.Walk(".", func(walkResult *fs.WalkResult) {
		if walkResult.Err != nil {
			glog.Fatalf("unexpected error: %v.", walkResult.Err)
		}
		sp.Add(&snapshot.FileItem{FullPath: walkResult.Path})
	})

	if err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}

	iter := sp.Iterator()
	for item := iter.Next(); item != nil; item = iter.Next() {
		fmt.Println(item.FullPath)
	}
}
