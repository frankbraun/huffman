/*

Package huffman implements a Huffman entropy coding.
https://en.wikipedia.org/wiki/Huffman_coding

*/
package huffman

import (
	"fmt"
	"sort"
	"strconv"
)

// Type of the value stored in a Node.
type ValueType int32

// A node in the Huffman tree.
type Node struct {
	Parent *Node     // Optional parent node, for fast code read-out
	Left   *Node     // Optional left node
	Right  *Node     // Optional right node
	Count  int       // Relative frequency
	Value  ValueType // Optional value, set if this is a leaf
}

// Code returns the Huffman code of the node.
// Left children get bit 0, Right children get bit 1.
// Implementation uses Node.Parent to walk "up" in the tree.
func (n *Node) Code() (r uint64, bits byte) {
	for parent := n.Parent; parent != nil; n, parent = parent, parent.Parent {
		if parent.Right == n { // bit 1
			r |= 1 << bits
		} // else bit 0 => nothing to do with r
		bits++
	}
	return
}

// Slice of *Node that implements sort.Interface, order defined by Node.Count.
type SortNodes []*Node

func (sn SortNodes) Len() int           { return len(sn) }
func (sn SortNodes) Less(i, j int) bool { return sn[i].Count < sn[j].Count }
func (sn SortNodes) Swap(i, j int)      { sn[i], sn[j] = sn[j], sn[i] }

// Build builds a Huffman tree from the specified leaves.
// The content of the passed slice is modified, if this is unwanted, pass a copy.
// Guaranteed that the same input slice will result in the same Huffman tree.
func Build(leaves []*Node) *Node {
	// We sort once and use binary insertion later on
	sort.Stable(SortNodes(leaves)) // Note: stable sort for deterministic output!

	return BuildSorted(leaves)
}

// BuildSorted builds a Huffman tree from the specified leaves which must be sorted by Node.Count.
// The content of the passed slice is modified, if this is unwanted, pass a copy.
// Guaranteed that the same input slice will result in the same Huffman tree.
func BuildSorted(leaves []*Node) *Node {
	if len(leaves) == 0 {
		return nil
	}

	for len(leaves) > 1 {
		left, right := leaves[0], leaves[1]
		parentCount := left.Count + right.Count
		parent := &Node{Left: left, Right: right, Count: parentCount}
		left.Parent = parent
		right.Parent = parent

		// Where to insert parent in order to remain sorted?
		ls := leaves[2:]
		idx := sort.Search(len(ls), func(i int) bool { return ls[i].Count >= parentCount })
		idx += 2

		// Insert
		copy(leaves[1:], leaves[2:idx])
		leaves[idx-1] = parent
		leaves = leaves[1:]
	}

	return leaves[0]
}

// Print traverses the Huffman tree and prints the values with their code in binary representation.
// For debugging purposes.
func Print(root *Node) {
	// traverse traverses a subtree from the given node,
	// using the prefix code leading to this node, having the number of bits specified.
	var traverse func(n *Node, code uint64, bits byte)

	traverse = func(n *Node, code uint64, bits byte) {
		if n.Left == nil {
			// Leaf
			fmt.Printf("'%c': %0"+strconv.Itoa(int(bits))+"b\n", n.Value, code)
			return
		}
		bits++
		traverse(n.Left, code<<1, bits)
		traverse(n.Right, code<<1+1, bits)
	}

	traverse(root, 0, 0)
}
