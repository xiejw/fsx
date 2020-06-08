package main

import (
	"fmt"

	"github.com/xiejw/lunar"
	"github.com/xiejw/lunar/scanner"

	"github.com/xiejw/fsx/go/snapshot"
)

func main() {
	lunar.Init(true /*parseFlag*/)
	defer lunar.FinishUp()

	sp := snapshot.New()

	scanner.Walk(".", []scanner.Filter{scanner.NewCommonFilter(nil)},
		func(metadata scanner.FileMetadata) {
			sp.Add(&snapshot.FileItem{FullPath: metadata.Path})
		})

	iter := sp.Iterator()
	for item := iter.Next(); item != nil; item = iter.Next() {
		fmt.Println(item.FullPath)
	}
}
