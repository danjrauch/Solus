package main

import (
	"errors"
	"math/bits"
	"sort"

	"github.com/golang-collections/collections/queue"
	"github.com/golang-collections/go-datastructures/bitarray"
	"github.com/hillbig/rsdic"
)

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

/*
	The work for these structures comes from the work of
	Brisaboa et al. Some of the paper titles are listed:

	"k2-trees for Compact Web Graph Representation"

	We can build a k-squared tree from adjacency lists by
	recursive descent using the theoretical structure below.

	Quadrant Ordering
	_____________
	|     |     |
	|  0  |  1  |
	|     |     |
	-------------
	|     |     |
	|  2  |  3  |
	|     |     |
	-------------

	We are input a list of adjacency lists that represent a graph.
	For each edge in the graph, we build a path in the k2-tree.
	Starting from the root, we insert k2-tree nodes based on
	the position of the edge in the graph's adjacency matrix.
	For example, if the edge in question lies in quadrant 2 of
	the adjacency matrix, we insert a k2-tree node into the
	children list for the root node if it doesn't exist already.
	Continue recursively into the found quadrant until the search
	space is one cell of the adjacency matrix.
*/

// K2Tree represents a k-squared tree
type K2Tree interface {
	GetChild(x int, c int) (bool, error)
}

type k2Tree struct {
	tree      bitarray.BitArray
	leaves    bitarray.BitArray
	lenTree   int
	lenLeaves int
	rank      *rsdic.RSDic
}

type k2TreeNode struct {
	children []*k2TreeNode
	value    bool
	level    int
}

// GetChild gets the cth child of node at pos x in tree
func (kt *k2Tree) GetChild(x int, c int) (bool, error) {
	n, err := kt.tree.GetBit(uint64(x))
	if err != nil {
		return false, err
	}

	if !n {
		return false, errors.New("Bit at pos x is not set")
	}

	pos := (int(kt.rank.Rank(uint64(x), true))+1)*4 + c

	if pos < kt.lenTree {
		return kt.tree.GetBit(uint64(pos))
	}

	return kt.leaves.GetBit(uint64(pos - kt.lenTree))
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
		root.value = true
		if root.children == nil {
			root.children = make([]*k2TreeNode, 4)
			for i := range root.children {
				root.children[i] = &k2TreeNode{value: false, level: root.level + 1}
			}
		}
		root.children[v].value = true
		root = root.children[v]
	}
}

func newK2Tree(graph [][]int) *k2Tree {
	nNodes := len(graph)
	root := k2TreeNode{value: true, level: 0}
	cursors := make([]int, nNodes)
	for _, row := range graph {
		sort.Ints(row)
	}
	for i := 0; i < nNodes; i++ {
		nEdges := len(graph[i])
		for j := 0; j < nNodes; j++ {
			if cursors[i] < nEdges {
				if graph[i][cursors[i]] == j {
					addK2TreeNode(&root, i, j, nextPowerOf2(uint(nNodes)))
					cursors[i]++
				}
			}
		}
	}

	maxLevel := 0

	qu := queue.New()
	qu.Enqueue(&root)
	for qu.Len() > 0 {
		var node *k2TreeNode = qu.Dequeue().(*k2TreeNode)
		if node.level > maxLevel {
			maxLevel = node.level
		}
		// fmt.Print(node.level, " ", node.value, " ")
		for _, child := range node.children {
			qu.Enqueue(child)
		}
	}
	// fmt.Println()

	var tree []bool
	var leaves []bool
	var rank *rsdic.RSDic = rsdic.New()

	qu.Enqueue(&root)
	for qu.Len() > 0 {
		var node *k2TreeNode = qu.Dequeue().(*k2TreeNode)
		if node.level != maxLevel && node.level != 0 {
			tree = append(tree, node.value)
			rank.PushBack(node.value)
		} else if node.level == maxLevel {
			leaves = append(leaves, node.value)
		}
		for _, child := range node.children {
			qu.Enqueue(child)
		}
	}

	// fmt.Println(tree)
	// fmt.Println(leaves)

	ktree := &k2Tree{
		tree:      bitarray.NewBitArray(uint64(len(tree))),
		leaves:    bitarray.NewBitArray(uint64(len(leaves))),
		lenTree:   len(tree),
		lenLeaves: len(leaves),
		rank:      rank,
	}

	for i, v := range tree {
		if v {
			ktree.tree.SetBit(uint64(i))
		}
	}

	for i, v := range leaves {
		if v {
			ktree.leaves.SetBit(uint64(i))
		}
	}

	return ktree
}

// NewK2Tree creates a new K2Tree
func NewK2Tree(graph [][]int) K2Tree {
	return newK2Tree(graph)
}
