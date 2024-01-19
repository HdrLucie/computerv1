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
	leftExpression := make(map[int]int)
	// rightExpression := make(map[int]int)

	fmt.Printf("%p\n%p\n", ast, ast.left.parent)
	fmt.Println(ast.left.parent.toStr())
	fmt.Println(ast.left.left.parent.parent.toStr())
	fmt.Println(ast.left.left.left.parent.toStr())
	fmt.Println(ast.left.left.left.parent.parent.parent.toStr())
	fmt.Println(ast.left.left.left.left.parent.parent.parent.parent.toStr())
	leftExpression = getExponentTab(ast, leftExpression)
	fmt.Println(leftExpression)
}

func getExponentTab(ast *Node, leftExpression map[int]int) map[int]int {
	// leftExpression := make(map[int]int)
	// rightExpression := make(map[int]int)

	if ast.parent != nil {
		fmt.Printf("Node : %s | MAMAN : %s | TATA : %s | MAMIE : %s\n", ast.toStr(), ast.parent.toStr(), ast.parent.sister().toStr(), ast.parent.parent.toStr())
	}
	if ast.node.tokenType == x {
		exponent := ast.sister()
		fmt.Printf("Node : %s | EXPO : %s | MAMAN : %s | TATA : %s | MAMIE : %s\n", ast.toStr(), exponent.toStr(), ast.parent.toStr(), ast.parent.sister().toStr(), ast.parent.parent.toStr())
		leftExpression[int(exponent.node.valueNumber)] = int(ast.parent.sister().node.valueNumber)
	}
	if ast.node.tokenType == digit {
		sister := ast.sister()
		fmt.Printf("Sister Value : %s | %s\n", sister.toStr(), sister.getType())
		if sister.node.tokenType == operator {
			fmt.Printf("value : %s\n", sister.toStr())
			leftExpression = getExponentTab(sister.left, leftExpression)
		}
	} else if ast.left != nil {
		fmt.Printf("VALUE : %s | PARENT : %s\n", ast.toStr(), ast.parent.toStr())
		leftExpression = getExponentTab(ast.left, leftExpression)
	}
	return leftExpression
}
