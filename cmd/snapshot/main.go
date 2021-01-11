package main

import (
	"fmt"

	"github.com/xiejw/fsx/src/snapshot"
)

func main() {

	r := snapshot.Region{
		Name:    "testing",
		BaseDir: ".",
	}

	err := snapshot.FillRegion(&r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("snapshot: \n")
	for _, item := range r.Items {
		fmt.Printf("  %10v: %v\n", item.Size, item.RelPath)
	}

}
