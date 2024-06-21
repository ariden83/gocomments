package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type FunctionInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <path/to/go/file>")
		return
	}

	filePath := os.Args[1]
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	var functions []FunctionInfo
	for _, f := range node.Decls {
		if fn, isFn := f.(*ast.FuncDecl); isFn {
			if fn.Doc != nil {
				functions = append(functions, FunctionInfo{
					Name:    fn.Name.Name,
					Comment: fn.Doc.Text(),
				})
			}
		}
	}

	jsonOutput, err := json.Marshal(functions)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonOutput))
}
