package main

import (
	"fmt"
)

// TreeNode represents a node in the binary tree
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// invertTree recursively inverts the binary tree
func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	// Swap left and right children
	root.Left, root.Right = root.Right, root.Left

	// Recursively invert left and right subtrees
	invertTree(root.Left)
	invertTree(root.Right)

	return root
}

// printTree prints the tree in a level-order traversal
func printTree(root *TreeNode) {
	if root == nil {
		fmt.Println("Empty tree")
		return
	}

	// Create a queue for level-order traversal
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		levelSize := len(queue)
		fmt.Print("Level: ")

		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]

			fmt.Printf("%d ", node.Val)

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		fmt.Println()
	}
}

func main() {
	// Create a sample binary tree:
	//       4
	//      / \
	//     2   7
	//    / \ / \
	//   1  3 6  9

	root := &TreeNode{
		Val: 4,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 1,
			},
			Right: &TreeNode{
				Val: 3,
			},
		},
		Right: &TreeNode{
			Val: 7,
			Left: &TreeNode{
				Val: 6,
			},
			Right: &TreeNode{
				Val: 9,
			},
		},
	}

	fmt.Println("Original Tree:")
	printTree(root)

	// Invert the tree
	invertedRoot := invertTree(root)

	fmt.Println("\nInverted Tree:")
	printTree(invertedRoot)
}
