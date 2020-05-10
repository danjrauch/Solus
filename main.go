package main

import (
	"fmt"
)

func main() {
	graph := [][]int{{0, 1}, {1}, {2, 3}, {2}}
	var ktree K2Tree = NewK2Tree(graph)
	val, err := ktree.GetChild(3, 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("(", 3, ",", 3, "):", val)
}
