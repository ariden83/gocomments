package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"

	"github.com/fatih/astrewrite"
)

type file struct {
	f        *ast.File
	fSet     *token.FileSet
	src      []byte
	fileName string
	cfg      *CommentConfig
}

func processComments(fileName string, src []byte, cache *CommentConfigCache) ([]byte, error) {
	fileSet := token.NewFileSet()

	if strings.HasSuffix(fileName, "_test.go") {
		return nil, nil
	} else if !strings.HasSuffix(fileName, ".go") {
		return nil, nil
	}

	f, err := parser.ParseFile(fileSet, fileName, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	cfg, err := cache.Get(fileName)
	if err != nil {
		return nil, err
	}

	file := file{
		cfg:      cfg,
		f:        f,
		src:      src,
		fileName: fileName,
		fSet:     fileSet,
	}

	return file.autoComment()
}

func (file *file) addSignature() string {
	if file.cfg.Signature != "" {
		return "\\ @author " + file.cfg.Signature + "."
	}
	return ""
}

func (file *file) autoComment() ([]byte, error) {
	var comments []*ast.CommentGroup

	for _, decl := range file.f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {

			switch genDecl.Tok {
			case token.TYPE:
				for _, spec := range genDecl.Specs {

					typeSpec := spec.(*ast.TypeSpec)

					if typeSpec.Doc.Text() == "" {
						privateValue := ""
						if !typeSpec.Name.IsExported() {
							privateValue = "private "
						}

						switch structType := typeSpec.Type.(type) {
						case *ast.StructType:
							txt := fmt.Sprintf("// %s represents a %sstructure for ", typeSpec.Name, privateValue)

							var mandatoryFields []*ast.Ident
							var optionFields []*ast.Ident

							// It contains information about the type of drop-off and
							// optionally, details about the sender's mailbox picking.
							for _, field := range structType.Fields.List {
								for _, f := range field.Names {
									if _, isPointer := field.Type.(*ast.StarExpr); isPointer {
										optionFields = append(optionFields, f)
									} else {
										mandatoryFields = append(mandatoryFields, f)
									}
								}
							}

							if len(mandatoryFields) > 0 {
								txt += "\n// It contains information about "
								for i, f := range mandatoryFields {
									if i > 0 {
										txt += ", "
									}
									privateKey := "private "
									if f.IsExported() {
										privateKey = ""
									}
									txt += fmt.Sprintf("%s %s%s", IndefiniteArticle(fmt.Sprintf("%s", f.Name)), privateKey, f.Name)
								}
							}

							if len(optionFields) > 0 {
								if len(mandatoryFields) > 0 {
									txt += " and\n// optionally "
								} else {
									txt += "\n// It contains optional information about "
								}

								for i, f := range optionFields {
									if i > 0 {
										txt += ", "
									}
									privateKey := "private "
									if f.IsExported() {
										privateKey = ""
									}
									txt += fmt.Sprintf("%s %s%s", IndefiniteArticle(fmt.Sprintf("%s", f.Name)), privateKey, f.Name)
								}
							}

							txt += file.addSignature()
							txt += ".\n"

							genDecl.Doc = &ast.CommentGroup{
								List: []*ast.Comment{{
									Text:  txt,
									Slash: typeSpec.Pos() - token.Pos(len("type ")+1),
								}},
							}

						default:
							txt := fmt.Sprintf("// %s is a type alias for the %s type.\n// It allows you to create a new type with the same\n// underlying type as int, but with a different name.\n// This can be useful for improving code readability\n// and providing more semantic meaning to your types.\n", typeSpec.Name, structType)
							txt += file.addSignature()
							genDecl.Doc = &ast.CommentGroup{
								List: []*ast.Comment{{
									Text:  txt,
									Slash: genDecl.Pos() - 1,
								}},
							}
						}
					}
				}
			case token.CONST:
				for _, spec := range genDecl.Specs {
					varSpec := spec.(*ast.ValueSpec)
					// Afficher le nom et le type de la variable

					hasParenthesis := false
					if genDecl.Lparen > 0 {
						hasParenthesis = true
					}

					for _, name := range varSpec.Names {
						if varSpec.Doc.Text() == "" {
							exported := ""
							if !name.IsExported() {
								exported = "private "
							}
							if decl, ok := name.Obj.Decl.(*ast.ValueSpec); ok {
								txt := fmt.Sprintf("// %s is a %sconstant which provides .", name.Name, exported)

								if hasParenthesis {
									decl.Doc = &ast.CommentGroup{
										List: []*ast.Comment{{
											Text:  txt,
											Slash: decl.Pos() - 1,
										}},
									}
								} else {
									txt += file.addSignature()
									genDecl.Doc = &ast.CommentGroup{
										List: []*ast.Comment{{
											Text:  txt,
											Slash: genDecl.Pos() - 1,
										}},
									}
								}
							}
						}
					}
				}
			case token.VAR:
				for _, spec := range genDecl.Specs {
					varSpec := spec.(*ast.ValueSpec)

					hasParenthesis := false
					if genDecl.Lparen > 0 {
						hasParenthesis = true
					}

					for _, name := range varSpec.Names {
						if varSpec.Doc.Text() == "" {
							exported := ""
							if !name.IsExported() {
								exported = "private "
							}

							if decl, ok := name.Obj.Decl.(*ast.ValueSpec); ok {
								txt := fmt.Sprintf("// %s is a %svariable of type %s which provides .", name.Name, exported, fmt.Sprintf("%s", decl.Type))

								if hasParenthesis {
									decl.Doc = &ast.CommentGroup{
										List: []*ast.Comment{{
											Text:  txt,
											Slash: decl.Pos() - 1,
										}},
									}

								} else {
									txt += file.addSignature()
									genDecl.Doc = &ast.CommentGroup{
										List: []*ast.Comment{{
											Text:  txt,
											Slash: genDecl.Pos() - 1,
										}},
									}
								}
							}
						}
					}
				}
			default:
			}
		}

		if genDecl, ok := decl.(*ast.FuncDecl); ok {
			if genDecl.Doc.Text() == "" && genDecl.Name.Name != "main" {
				txt := getFuncComments(genDecl)
				txt += file.addSignature()
				genDecl.Doc = &ast.CommentGroup{
					List: []*ast.Comment{{
						Text:  txt,
						Slash: genDecl.Pos() - 1,
					}},
				}
			}
		}
	}

	file.f.Comments = comments

	reWriteFunc := func(node ast.Node) (ast.Node, bool) {
		return node, true
	}

	newAst := astrewrite.Walk(file.f, reWriteFunc)
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, file.fSet, newAst); err != nil {
		log.Fatal(err)
	}
	log.Printf("result %s", buf.String())
	return buf.Bytes(), nil
}

func isNewFunc(name string) bool {
	return strings.HasPrefix(name, "New")
}

func newFuncTxt(fn *ast.FuncDecl) string {
	instanceReturnMsg := ""
	initializesMsg := ""
	errorReturnMsg := ""
	if fn.Type.Results != nil {
		if len(fn.Type.Results.List) != 0 {
			for i, res := range fn.Type.Results.List {
				typeReturnKey := fmt.Sprintf("%s", res.Type)
				if typeReturnKey != "error" {
					if i > 0 {
						instanceReturnMsg += " and "
					} else if i == 0 {
						instanceReturnMsg += " of "
						initializesMsg = typeReturnKey
					}
					instanceReturnMsg += typeReturnKey

				} else {
					errorReturnMsg = "\n// It's return an error if the initialization fails, otherwise nil."
				}
			}
		}
	}

	txt := fmt.Sprintf("// %s creates a new instance%s.", fn.Name.Name, instanceReturnMsg)
	if (fn.Type.Params == nil || len(fn.Type.Params.List) == 0) && (fn.Type.Results == nil || len(fn.Type.Results.List) == 0) {
	} else {
		if fn.Type.Params != nil {
			if len(fn.Type.Params.List) != 0 {
				txt += fmt.Sprintf("\n// It initializes the %s with the provided ", initializesMsg)
			}

			for i, param := range fn.Type.Params.List {
				if i > 0 {
					txt += ", "
				}
				for _, name := range param.Names {
					txt += fmt.Sprintf("%s of type %s", name.Name, getTypeName(param.Type))
				}
			}
		}
		txt += "."
	}
	if errorReturnMsg != "" {
		txt += errorReturnMsg
	}

	return txt
}

func getFuncComments(fn *ast.FuncDecl) string {
	txt := ""
	if isNewFunc(fn.Name.Name) {
		return newFuncTxt(fn)
	}

	privateValue := ""
	if !fn.Name.IsExported() {
		privateValue = "private "
	}
	if (fn.Type.Params == nil || len(fn.Type.Params.List) == 0) && (fn.Type.Results == nil || len(fn.Type.Results.List) == 0) {
		if fn.Recv != nil {
			funcType := ""
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				if t, ok := fn.Recv.List[0].Type.(*ast.Ident); ok {
					funcType = t.String()
				}
			}
			txt = fmt.Sprintf("// %s is a %smethod that belongs to the %s struct.\n// It does not take any arguments.\n", fn.Name.Name, privateValue, funcType)
		} else {
			txt = fmt.Sprintf("// %s is a %smethod .\n// It does not take any arguments.\n", fn.Name.Name, privateValue)
		}

	} else {
		if fn.Recv != nil {
			funcType := ""
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				if t, ok := fn.Recv.List[0].Type.(*ast.Ident); ok {
					funcType = t.String()
				}
			}
			txt = fmt.Sprintf("// %s is a %smethod that belongs to the %s struct", fn.Name.Name, privateValue, funcType)
		} else {
			txt = fmt.Sprintf("// %s is a %smethod", fn.Name.Name, privateValue)
		}

		if fn.Type.Params != nil {
			if len(fn.Type.Params.List) != 0 {
				txt += " that take "
			}

			for i, param := range fn.Type.Params.List {
				if i > 0 {
					txt += ", "
				}
				for _, name := range param.Names {
					txt += fmt.Sprintf("%s %s of type %s", IndefiniteArticle(fmt.Sprintf("%s", name.Name)), name.Name, getTypeName(param.Type))
				}
			}
		}

		if fn.Type.Results != nil && len(fn.Type.Results.List) != 0 {
			errorReturnMsg := ""
			for i, res := range fn.Type.Results.List {
				typeReturnKey := fmt.Sprintf("%s", res.Type)
				if typeReturnKey == "error" {
					errorReturnMsg = "\n// It's return an error if fails, otherwise nil"
				} else {
					if i > 0 {
						txt += " and "
					} else {
						txt += "\n// and returns "
					}
					txt += fmt.Sprintf("%s %s", IndefiniteArticle(fmt.Sprintf("%s", res.Type)), res.Type)
				}
			}
			txt += errorReturnMsg
		}
		txt += ".\n"
	}
	return txt
}

func IndefiniteArticle(word string) string {
	// Convertir le mot en minuscules pour faciliter la comparaison.
	word = strings.ToLower(word)
	// Les lettres qui nécessitent "an".
	anLetters := "aeiou"
	// Si le mot commence par une voyelle, retourner "an".
	if strings.ContainsRune(anLetters, rune(word[0])) {
		return "an"
	}
	// Sinon, retourner "a".
	return "a"
}

/*
func getTypeNameForVar(expr any) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.AssignStmt:
		return "AssignStmt"
	case *ast.StarExpr:
		return "*" + getTypeName(expr.X)
	case *ast.SelectorExpr:
		pkg := getTypeName(expr.X)
		sel := expr.Sel.Name
		return pkg + "." + sel
	case *ast.ArrayType:
		return "[]" + getTypeName(expr.Elt)
	case *ast.MapType:
		return "map[" + getTypeName(expr.Key) + "]" + getTypeName(expr.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		dir := ""
		switch expr.Dir {
		case ast.RECV:
			dir = "<-chan "
		case ast.SEND:
			dir = "chan<- "
		}
		return dir + getTypeName(expr.Value)
	case *ast.FuncType:
		// Handle function types if needed
		// You can recursively call getTypeName for Params and Results
		return "func"
	default:
		return "unknown"
	}
}*/

func getTypeName(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		return "*" + getTypeName(expr.X)
	case *ast.SelectorExpr:
		pkg := getTypeName(expr.X)
		sel := expr.Sel.Name
		return pkg + "." + sel
	case *ast.ArrayType:
		return "[]" + getTypeName(expr.Elt)
	case *ast.MapType:
		return "map[" + getTypeName(expr.Key) + "]" + getTypeName(expr.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		dir := ""
		switch expr.Dir {
		case ast.RECV:
			dir = "<-chan "
		case ast.SEND:
			dir = "chan<- "
		}
		return dir + getTypeName(expr.Value)
	case *ast.FuncType:
		// Handle function types if needed
		// You can recursively call getTypeName for Params and Results
		return "func"
	default:
		return "unknown"
	}
}
