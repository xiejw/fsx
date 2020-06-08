package main

import (
	"fmt"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/xiejw/lunar"
	"github.com/xiejw/lunar/scanner"

	"github.com/xiejw/fsx/go/snapshot"
)

func main() {
	lunar.Init(true /*parseFlag*/)
	defer lunar.FinishUp()

	sp := snapshot.New()

	path, err := filepath.Abs(".")
	if err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}

	scanner.Walk(path, []scanner.Filter{scanner.NewCommonFilter([]string{})},
		func(metadata scanner.FileMetadata) {
			info := metadata.Info
			if info.IsDir() {
				return
			}
			sp.Add(&snapshot.FileItem{FullPath: metadata.Path})
		})

	iter := sp.Iterator()
	for item := iter.Next(); item != nil; item = iter.Next() {
		fmt.Println(item.FullPath)
	}
}
