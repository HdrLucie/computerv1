package main

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

type Node struct {
	node   Token
	parent *Node
	left   *Node
	right  *Node
}

type Tree struct {
	head *Node
}

func (self *Node) toStr() string {
	if self == nil {
		return "nil"
	}
	if self.node.valueString == "" {
		return fmt.Sprintf("%f", self.node.valueNumber)
	}
	return self.node.valueString
}

func (self *Node) getType() enumType {
	return self.node.tokenType
}

func (self *Node) sister() *Node {
	if self.parent == nil {
		return nil
	}
	fmt.Printf("Parent ouhouhou: %s\n", self.parent.toStr())
	if self.parent.left == self {
		return self.parent.right
	}
	return self.parent.left
}

type parseError struct{}

func (m *parseError) Error() string {
	return "parse error (wrong input)"
}

func newNode(token Token) Node {
	newNode := Node{node: token, parent: nil, left: nil, right: nil}
	return newNode
}

func tokenToNode(tokenTab []Token) []*Node {
	var nodeTab []*Node

	for _, token := range tokenTab {
		tmpNode := newNode(token)
		nodeTab = append(nodeTab, &tmpNode)
	}
	return nodeTab
}

func linkNode(current *Node, left *Node, right *Node) *Node {
	current.left = left
	current.right = right
	left.parent = current
	right.parent = current

	return current
}

func reduceNodeTab(nodeTab []*Node, isOperator func(string) bool) []*Node {

	modified := true
	for modified {
		modified = false
		var indexes []int
		for i, node := range nodeTab {
			if isOperator(node.node.valueString) && node.left == nil && node.right == nil && i != 0 && i != len(nodeTab)-1 {
				indexes = append(indexes, i-1)
				indexes = append(indexes, i+1)
				nodeTab[i] = linkNode(nodeTab[i], nodeTab[i-1], nodeTab[i+1])
				modified = true
				break
			}
		}
		var newNodeTab []*Node
		for i := range nodeTab {
			if slices.Contains(indexes, i) {
				continue
			}
			newNodeTab = append(newNodeTab, nodeTab[i])
		}
		nodeTab = newNodeTab
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

func createAST(tokenTab []Token) (*Node, error) {
	nodeTab := tokenToNode(tokenTab)

	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "^" })
	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "*" || s == "/" })
	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "+" || s == "-" })
	nodeTab = reduceNodeTab(nodeTab, func(s string) bool { return s == "=" })

	if len(nodeTab) != 1 {
		return nodeTab[0], &parseError{}
	}
	return nodeTab[0], nil
}

func nodeStr(node *Node) string {
	if node.node.valueString != "" {
		return node.node.valueString
	}
	return fmt.Sprintf("%.9g", node.node.valueNumber)
}

func print(file *os.File, node *Node) {
	// fmt.Print("-> ")
	// fmt.Println(*node)
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
	// fmt.Println("Creating dotgraph")
	// fmt.Println("root is:", node)
	file, errs := os.Create("graph.dot")
	if errs != nil {
		// fmt.Println("Failed to create file:", errs)
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
