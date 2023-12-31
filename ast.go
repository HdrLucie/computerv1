package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type Node struct {
	node  Token
	left  *Node
	right *Node
}

type Tree struct {
	head *Node
}

type parseError struct{}

func (m *parseError) Error() string {
	return "parse error (wrong input)"
}

func newNode(token Token) Node {
	newNode := Node{node: token, left: nil, right: nil}
	return newNode
}

func tokenToNode(tokenTab []Token) []Node {
	var nodeTab []Node

	for _, token := range tokenTab {
		nodeTab = append(nodeTab, newNode(token))
	}
	return nodeTab
}

func linkNode(parent Node, left Node, right Node) Node {
	parent.left = &left
	parent.right = &right

	return parent
}

func reduceNodeTab(nodeTab []Node, isOperator func(string) bool) []Node {
	modified := true
	for modified {
		modified = false
		var indexes []int
		for i, node := range nodeTab {
			if isOperator(node.node.valueString) && node.left == nil && node.right == nil && i != 0 && i != len(nodeTab)-1 {
				indexes = append(indexes, i-1)
				indexes = append(indexes, i+1)
				nodeTab[i] = linkNode(nodeTab[i], nodeTab[i-1], nodeTab[i+1])
				//fmt.Println(nodeTab[i], nodeTab[i].left, nodeTab[i].right)
				modified = true
				break
			}
		}
		var newNodeTab []Node
		for i := range nodeTab {
			if slices.Contains(indexes, i) {
				continue
			}
			newNodeTab = append(newNodeTab, nodeTab[i])
		}
		nodeTab = newNodeTab
		fmt.Println("> ", nodeTab)
	}
	return nodeTab
}

func printOperation(node *Node) {
	if node.left != nil {
		printOperation(node.left)
	}
	fmt.Print(node.node)
	if node.right != nil {
		printOperation(node.right)
	}
}

func createAST(tokenTab []Token) (Node, error) {
	nodeTab := tokenToNode(tokenTab)

	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "^" })
	fmt.Println(nodeTab)
	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "*" || s == "/" })
	fmt.Println(nodeTab)
	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "+" || s == "-" })
	fmt.Println(nodeTab)
	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "=" })

	fmt.Println(nodeTab)
	// printNode(&nodeTab[0])
	if len(nodeTab) != 1 {
		return nodeTab[0], &parseError{}
	}
	// printOperation(&nodeTab[0])
	// fmt.Println()
	// fmt.Println(nodeTab)
	return nodeTab[0], nil
}

func nodeStr(node *Node) string {
	if node.node.valueString != "" {
		return node.node.valueString
	}
	return fmt.Sprintf("%.9g", node.node.valueNumber)
}

func print(file *os.File, node *Node) {
	fmt.Print("-> ")
	fmt.Println(*node)
	var out string
	out += "\t" + fmt.Sprintf("%#p", node) + " [label=\"" + nodeStr(node) + "\"]\n"
	// out += ";\n"
	if node.left != nil {
		print(file, node.left)
		out += "\t" + fmt.Sprintf("%#p", node) + " -> " + fmt.Sprintf("%#p", node.left) + "[label=\"L\"];\n"
	}
	if node.right != nil {
		print(file, node.right)
		out += "\t" + fmt.Sprintf("%#p", node) + " -> " + fmt.Sprintf("%#p", node.right) + "[label=\"R\"];\n"
	}
	file.WriteString(out)
}

func printNode(node *Node) {
	fmt.Println("Creating dotgraph")
	fmt.Println("root is:", node)
	file, errs := os.Create("graph.dot")
	if errs != nil {
		fmt.Println("Failed to create file:", errs)
		return
	}
	_, errs = file.WriteString("digraph {\n\tedge [arrowhead=none];\n\n")
	// single node with label saying empty
	if node == nil {
		_, errs = file.WriteString("\tempty [label=\"Rien.\", shape=none];\n")
	} else {
		print(file, node)
	}
	_, errs = file.WriteString("}\n")
	defer file.Close()
	return
}
