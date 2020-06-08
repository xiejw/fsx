package main

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/xiejw/lunar"
	"github.com/xiejw/lunar/scanner"

	"github.com/xiejw/fsx/go/snapshot"
)

func main() {
	lunar.Init(true /*parseFlag*/)
	defer lunar.FinishUp()

	sp := snapshot.New()

	var err error
	err = scanner.Walk(".",
		[]scanner.Filter{scanner.NewCommonFilter([]string{})},
		func(metadata scanner.FileMetadata) {
			info := metadata.Info
			if info.IsDir() {
				return
			}
			sp.Add(&snapshot.FileItem{FullPath: metadata.Path})
		})

	if err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}

	iter := sp.Iterator()
	for item := iter.Next(); item != nil; item = iter.Next() {
		fmt.Println(item.FullPath)
	}
}
