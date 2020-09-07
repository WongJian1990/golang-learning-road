package main

import (
	"crypto/sha256"
)

//MerkleTree MerkleTree
type MerkleTree struct {
	RootNode *MerkleNode
}

//MerkleNode MerkleTree节点
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

//NewMerkleNode 构造MerkleNode
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}
	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {

		leaveHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(leaveHash)
		node.Data = hash[:]
	}
	node.Left = left
	node.Right = right
	return node
}

//NewMerkleTree 构造Merkle树
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode
	length := len(data)
	if length%2 != 0 {
		data = append(data, data[length-1])
	}

	for _, v := range data {
		node := NewMerkleNode(nil, nil, v)
		nodes = append(nodes, *node)
	}

	for len(nodes) > 1 {
		var newLevel []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}
		nodes = newLevel
	}
	mTree := MerkleTree{&nodes[0]}
	return &mTree
}
