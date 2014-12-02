package main

import (
	"fmt"
	"github.com/Stantheman/pstree"
)

func main() {
	tree := make(pstree.ProcessTree)

	if err := tree.Populate(); err != nil {
		fmt.Errorf("Failed getting proceses: %v\n", err)
		return
	}

	fmt.Print(tree)
}
