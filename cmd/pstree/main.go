package main

import (
	"fmt"
	"github.com/Stantheman/pstree"
)

const IndentDepth = 0

func main() {
	tree := make(pstree.ProcessTree)

	if err := tree.Populate(); err != nil {
		fmt.Errorf("Failed getting proceses: %v\n", err)
		return
	}

	tree.PrintDepthFirst("0", IndentDepth)
}
