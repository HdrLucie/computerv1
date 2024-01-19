package main

import (
	"fmt"
	"os"
)

func main() {
	tokensTab, err := lexerFunction(os.Args[1])
	if err != nil {
		fmt.Printf("Error - %s\n", err)
		os.Exit(1)
	}
	// fmt.Println(tokensTab)
	ast, err := createAST(tokensTab)
	if err != nil {
		fmt.Printf("Error - %s\n", err)
		os.Exit(1)
	}
	printNode(ast)
	createExponentTab(ast)
}
