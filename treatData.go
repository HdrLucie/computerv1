package main

import "fmt"

// "fmt"

func getLastNode(root *Node) *Node {
	var end *Node = root
	for end.left != nil {
		end = end.left
	}
	return end
}

func createExponentTab(ast *Node) {
	// leftExpression := make(map[int]int)
	// rightExpression := make(map[int]int)

	if ast.node.tokenType == x {
		sister := ast.sister()

	}
	if ast.node.tokenType == digit {
		fmt.Printf("In sister Value : %s | %s\n", ast.toStr(), ast.getType())
		sister := ast.sister()
		fmt.Printf("Sister Value : %s | %s\n", sister.toStr(), sister.getType())
		if sister.node.tokenType == operator {
			fmt.Printf("value : %s\n", sister.toStr())
			createExponentTab(sister.left)
		}
	} else if ast.left != nil {
		fmt.Printf("VALUE : %s | PARENT : %s\n", ast.toStr(), ast.parent.toStr())
		createExponentTab(ast.left)
	}
}
