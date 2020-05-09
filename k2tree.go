package main

import (
	"fmt"
	"math/bits"
	"sort"

	"github.com/golang-collections/collections/queue"
	"github.com/golang-collections/go-datastructures/bitarray"
)

type K2Tree interface {
}

type k2Tree struct {
	tree   bitarray.BitArray
	leaves bitarray.BitArray
}

type k2TreeNode struct {
	children []*k2TreeNode
	value    bool
}

func prevPowerOf2(n uint) int {
	var u uint = 0
	for n > 0 {
		n >>= 1
		u++
	}
	return 1 << (u - 1)
}

func nextPowerOf2(n uint) int {
	if bits.OnesCount(n) == 1 {
		return int(n)
	}
	var u uint = 0
	for n > 0 {
		n >>= 1
		u++
	}
	return 1 << u
}

func addK2TreeNode(root *k2TreeNode, row int, col int, n int) {
	var path []int
	k := 2
	for n/k >= 1 {
		blockSize := n / k
		if row < blockSize && col < blockSize {
			// In quadrant 0
			path = append(path, 0)
		} else if row < blockSize && col >= blockSize {
			// In quadrant 1
			path = append(path, 1)
			col -= blockSize
		} else if row >= blockSize && col < blockSize {
			// In quadrant 2
			path = append(path, 2)
			row -= blockSize
		} else {
			// In quadrant 3
			path = append(path, 3)
			row -= blockSize
			col -= blockSize
		}
		k *= 2
	}
	for _, v := range path {
		if root.value == false || root.children == nil {
			root.value = true
			root.children = make([]*k2TreeNode, 4)
			for i := range root.children {
				root.children[i] = &k2TreeNode{value: false}
			}
		}
		root.children[v].value = true
		root = root.children[v]
	}
}

func newK2Tree(graph [][]int) *k2Tree {
	nNodes := len(graph)
	root := k2TreeNode{value: true}
	cursors := make([]int, nNodes)
	for _, row := range graph {
		sort.Ints(row)
	}
	for i := 0; i < nNodes; i++ {
		nEdges := len(graph[i])
		for j := 0; j < nNodes; j++ {
			if cursors[i] < nEdges {
				if graph[i][cursors[i]] == j {
					// TODO build node in tree
					addK2TreeNode(&root, i, j, nextPowerOf2(uint(nNodes)))
					cursors[i]++
				}
			}
		}
	}

	qu := queue.New()
	qu.Enqueue(&root)
	for qu.Len() > 0 {
		var node *k2TreeNode = qu.Dequeue().(*k2TreeNode)
		fmt.Print(node.value, " ")
		for _, child := range node.children {
			qu.Enqueue(child)
		}
	}
	fmt.Println()

	return &k2Tree{}
}

// NewK2Tree creates a new K2Tree
func NewK2Tree(graph [][]int) K2Tree {
	return newK2Tree(graph)
}
