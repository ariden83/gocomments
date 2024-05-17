package pkg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// Code source Go à analyser
	src := `
package main

const (
	ConstVar1 = 10
	ConstVar2 = "Hello"
)

const ConstVar3 = 10

var (
	Var1 int
	Var2 string
)
`

	// Analyse du code source
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		fmt.Println("Erreur lors de l'analyse du code source:", err)
		return
	}

	// Parcourir les déclarations dans le fichier analysé
	for _, decl := range node.Decls {
		// Vérifier si la déclaration est une déclaration générale (variables, constantes ou types)
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Vérifier le type de déclaration (const, var, type)
			switch genDecl.Tok {
			case token.CONST:
				fmt.Println("Déclaration de constante trouvée:")
				// Parcourir les spécifications de déclaration pour obtenir les détails des constantes
				for _, spec := range genDecl.Specs {
					constSpec := spec.(*ast.ValueSpec)
					// Afficher le nom et la valeur de la constante
					fmt.Printf("Nom: %s, Valeur: %s\n", constSpec.Names[0], constSpec.Values[0])
				}
			case token.VAR:
				fmt.Println("Déclaration de variable trouvée:")
				// Parcourir les spécifications de déclaration pour obtenir les détails des variables
				for _, spec := range genDecl.Specs {
					varSpec := spec.(*ast.ValueSpec)
					// Afficher le nom et le type de la variable
					fmt.Printf("Nom: %s, Type: %s\n", varSpec.Names[0], varSpec.Type)
				}
			}
		}
	}
}
