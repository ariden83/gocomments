package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/ariden/gocomments/internal/comments"
)

type FunctionInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// Mots-clés à ignorer dans les commentaires
var keywordsToIgnore = []string{"fixme", "note", "go:embed", "todo", "bug", "deprecated"}

// Fonction pour vérifier si un commentaire contient l'un des mots-clés interdits.
func containsIgnoredKeywords(comment string) bool {
	lowerComment := strings.ToLower(comment)
	for _, keyword := range keywordsToIgnore {
		if strings.Contains(lowerComment, keyword) {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		// fmt.Println("Usage: go run script.go <path/to/go/file>")
		return
	}

	args := os.Args
	var goFiles []string
	var otherArgs []string

	// Parse the arguments
	for i, arg := range args {
		if arg == "--" {
			otherArgs = args[i+1:]
			break
		}
		goFiles = append(goFiles, arg)
	}

	filePath := otherArgs[0]
	fSet := token.NewFileSet()

	node, err := parser.ParseFile(fSet, filePath, nil, parser.ParseComments)
	if err != nil {
		// fmt.Println(err)
		return
	}

	var functions []FunctionInfo
	for _, f := range node.Decls {
		if fn, isFn := f.(*ast.FuncDecl); isFn {
			if fn.Doc != nil && fn.Doc.Text() != "" && fn.Name.Name != "main" && fn.Name.Name != "init" {
				comment := fn.Doc.Text()
				if !containsIgnoredKeywords(comment) {
					functions = append(functions, FunctionInfo{
						Name:    comments.GenerateFuncCode(fn),
						Comment: comment,
					})
				}
			}
		}
	}

	jsonOutput, err := json.Marshal(functions)
	if err != nil {
		// fmt.Println(nil)
		return
	}

	fmt.Println(string(jsonOutput))
}
