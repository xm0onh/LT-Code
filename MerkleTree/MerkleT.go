package MerkleTree

import (
	"crypto/sha256"
	"fmt"
)

type TreeNode struct {
	Value []byte
	//Parent []byte
	LChild []byte
	RChild []byte
	Level  int
	IsRoot bool
	IsLeaf bool
}

type MerkleTree struct {
	Node []TreeNode
	//Root []byte
}

func NewTree(HashSlice [][]byte) *MerkleTree {
	//	LeavesCreated := false
	//if LeavesCreated == false {
	merkleTree := CreateLeaves(HashSlice)
	//LeavesCreated = true
	//	}
	parentSlice := make([][]byte, 0, 20)
	for len(HashSlice) > 1 {
		fmt.Println("Len HashSlice outside the second Loop is", len(HashSlice))

		for i := 0; i < len(HashSlice); i = i + 2 {
			fmt.Println("Len HashSlice inside the second Loop is", len(HashSlice))
			parent := make([]byte, 0, 10)
			ntree := CreateTreeNode()
			//		fmt.Println("The value of i is", i)
			//	fmt.Println("HashSlice[i]", HashSlice[i])
			//	fmt.Println("Len HashSlice[i]", len(HashSlice))

			ntree.LChild = HashSlice[i]
			parent = append(parent, HashSlice[i]...)
			if i+1 < len(HashSlice) {
				parent = append(parent, HashSlice[i+1]...)
				ntree.RChild = HashSlice[i+1]

			}
			parentHash, err := CalcHash(parent)
			if err != nil {
				fmt.Println("Error while generating Hahs ", err)
			}
			ntree.Value = parentHash
			//	merkleTree.Node = append(merkleTree.Node, *ntree)
			//	fmt.Println("ParentHash is ",parentHash )
			parentSlice = append(parentSlice, parentHash)
			ntree.Value = parentHash
			merkleTree.Node = append(merkleTree.Node, *ntree)
			//HashSlice = parentSlice
			//fmt.Println("Parent slice is", HashSlice)
			//	fmt.Println("ntree is", ntree)

		}
		HashSlice = parentSlice
		fmt.Println("Len Parent slice is", len(parentSlice))
		parentSlice = parentSlice[:0]

		//	if len(HashSlice)==1{
		//		return merkleTree
		//	}

	}

	fmt.Println("Merkle Tree is ", merkleTree)
	return merkleTree

}

func CreateLeaves(HashSlice [][]byte) *MerkleTree {
	merkleTree := CreateEmptyMerkleTree()
	for _, value := range HashSlice {
		treeNode := CreateTreeNode()
		treeNode.IsLeaf = true
		treeNode.Value = value
		merkleTree.Node = append(merkleTree.Node, *treeNode)
	}
	return merkleTree
}

func CreateTreeNode() *TreeNode {
	return &TreeNode{
		Value: nil,
		//	Parent: nil,
		LChild: nil,
		RChild: nil,
		Level:  0,
		IsRoot: false,
		IsLeaf: false,
	}
}

func CreateEmptyMerkleTree() *MerkleTree {
	return &MerkleTree{
		Node: make([]TreeNode, 0, 20),
		//Root: nil,
	}

}

//TestSHA256Content implements the Content interface provided by merkletree and represents the content stored in the tree.
type TestSHA256Content struct {
	x string
}

//CalculateHash hashes the values of a TestSHA256Content
func CalcHash(ChildrenConcatHash []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(ChildrenConcatHash); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
