package main

import (
	"bytes"

	ssz "github.com/ferranbt/fastssz"
)

// graftSubtree finds the node with the given value in the tree and replaces the parent's pointer to it with newNode.
func graftSubtree(root *ssz.Node, value []byte, newNode *ssz.Node) bool {
	// If the root is nil, the tree is empty or we've reached a leaf, return false
	if root == nil {
		return false
	}

	// Check if the left child matches the target value
	if root.Left() != nil && bytes.Equal(root.Left().Value(), value) {
		// Replace the left child with newNode
		root.SetLeft(newNode)
		return true
	}

	// Check if the right child matches the target value
	if root.Right() != nil && bytes.Equal(root.Right().Value(), value) {
		// Replace the right child with newNode
		root.SetRight(newNode)
		return true
	}

	// Recursively search in the left subtree
	if graftSubtree(root.Left(), value, newNode) {
		return true
	}

	// Recursively search in the right subtree
	if graftSubtree(root.Right(), value, newNode) {
		return true
	}

	// Value not found
	return false
}
