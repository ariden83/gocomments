package comments

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
	"github.com/stoewer/go-strcase"
)

type file struct {
	f         *ast.File
	fSet      *token.FileSet
	src       []byte
	fileName  string
	cfg       *CommentConfig
	processor commentsProcess
}

func Process(fileName string, src []byte, cache *CommentConfigCache) ([]byte, error) {
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

	processor := newProcessor(cfg)

	file := file{
		cfg:       cfg,
		processor: processor,
		f:         f,
		src:       src,
		fileName:  fileName,
		fSet:      fileSet,
	}

	return file.autoComment()
}

func (file *file) autoComment() ([]byte, error) {
	var comments []*ast.CommentGroup

	for _, decl := range file.f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {

			switch genDecl.Tok {
			case token.TYPE:
				if err := file.commentType(genDecl); err != nil {
					return nil, err
				}
			case token.CONST:
				if err := file.commentConst(genDecl); err != nil {
					return nil, err
				}
			case token.VAR:
				if err := file.commentVar(genDecl); err != nil {
					return nil, err
				}
			default:
			}
		}

		if genDecl, ok := decl.(*ast.FuncDecl); ok {
			if err := file.commentFunc(genDecl); err != nil {
				return nil, err
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

	return buf.Bytes(), nil
}

func GenerateFuncCode(fn *ast.FuncDecl) string {
	functionName := fn.Name.Name

	var (
		input    string
		output   string
		funcType string
	)

	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		for i, param := range fn.Type.Params.List {
			if i > 0 {
				input += ", "
			}
			for _, name := range param.Names {
				input += fmt.Sprintf("%s %s", name.Name, getTypeName(param.Type))
			}
		}
	}

	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		for i, res := range fn.Type.Results.List {
			if i > 0 {
				output += ", "
			}
			output += fmt.Sprintf("%s", getTypeName(res.Type))
		}
	}

	if fn.Recv != nil {
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			funcTypeName := ""
			txt := ""
			switch expr := fn.Recv.List[0].Type.(type) {
			case *ast.Ident:
				funcTypeName = expr.Name
				txt = funcTypeName
			case *ast.StarExpr:
				if f, ok := expr.X.(*ast.Ident); ok {
					txt = "*" + f.Name
					funcTypeName = f.Name
				}
			}

			funcType = "(" + strings.ToLower(string(funcTypeName[0])) + " " + txt + ") "
		}
	}

	return fmt.Sprintf(`func %s%s(%s) (%s)`, funcType, functionName, input, output)
}

func (file *file) commentConst(genDecl *ast.GenDecl) error {
	for _, spec := range genDecl.Specs {
		varSpec := spec.(*ast.ValueSpec)
		hasParenthesis := false
		if genDecl.Lparen > 0 {
			hasParenthesis = true
		}

		for _, name := range varSpec.Names {
			if varSpec.Doc.Text() == "" {
				exported := true
				if !name.IsExported() {
					exported = false
				}
				if decl, ok := name.Obj.Decl.(*ast.ValueSpec); ok {
					txt, err := file.processor.commentConst(name.Name, exported)
					if err != nil {
						return fmt.Errorf("fail to add comments on const: %v", err)
					}

					if hasParenthesis {
						decl.Doc = &ast.CommentGroup{
							List: []*ast.Comment{{
								Text:  txt,
								Slash: decl.Pos() - 1,
							}},
						}
					} else {
						txt += "\n"
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
	return nil
}

func (file *file) addSignature() string {
	if file.cfg.Signature != nil && *file.cfg.Signature != "" {
		return "//\n// Author: " + *file.cfg.Signature + "."
	}
	return ""
}

func (file *file) commentVar(genDecl *ast.GenDecl) error {
	for _, spec := range genDecl.Specs {
		varSpec := spec.(*ast.ValueSpec)

		hasParenthesis := false
		if genDecl.Lparen > 0 {
			hasParenthesis = true
		}

		for _, name := range varSpec.Names {
			if varSpec.Doc.Text() == "" {
				exported := true
				if !name.IsExported() {
					exported = false
				}

				explainVar := convertVarToCamelCaseTo(name.Name)

				if decl, ok := name.Obj.Decl.(*ast.ValueSpec); ok {

					txt, err := file.processor.commentVar(name.Name, fmt.Sprintf("%s", decl.Type), explainVar, exported)
					if err != nil {
						return fmt.Errorf("fail to add comments on var: %v", err)
					}

					if hasParenthesis {
						decl.Doc = &ast.CommentGroup{
							List: []*ast.Comment{{
								Text:  txt,
								Slash: decl.Pos() - 1,
							}},
						}

					} else {
						txt += ".\n"
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
	return nil
}

func (file *file) commentType(genDecl *ast.GenDecl) error {
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
						txt += fmt.Sprintf("%s %s%s", indefiniteArticle(fmt.Sprintf("%s", f.Name)), privateKey, f.Name)
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
						txt += fmt.Sprintf("%s %s%s", indefiniteArticle(fmt.Sprintf("%s", f.Name)), privateKey, f.Name)
					}
				}

				txt += ".\n"
				txt += file.addSignature()

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

	return nil
}

func isNewFunc(name string) bool {
	return strings.HasPrefix(name, "New")
}

// RequestPayload defines the structure of the request payload to the Anthropic API
type RequestPayload struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
}

// ResponsePayload defines the structure of the response payload from the Anthropic API
type ResponsePayload struct {
	Completion string `json:"completion"`
}

func (file *file) commentFunc(genDecl *ast.FuncDecl) error {
	if genDecl.Doc.Text() == "" && genDecl.Name.Name != "main" && genDecl.Name.Name != "init" {
		txt, err := file.processor.commentFunc(genDecl)
		if err != nil {
			log.Printf("fail to generate comment for func %s: %+v", genDecl.Name.Name, err)
			return err
		}
		txt += file.addSignature()
		genDecl.Doc = &ast.CommentGroup{
			List: []*ast.Comment{{
				Text:  txt,
				Slash: genDecl.Pos() - 1,
			}},
		}
	}
	return nil
}

func indefiniteArticle(word string) string {
	word = strings.ToLower(word)
	anLetters := "aeiou"
	if strings.ContainsRune(anLetters, rune(word[0])) {
		return "an"
	}
	return "a"
}

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

func convertVarToCamelCaseTo(str string) string {
	txt := strcase.SnakeCase(str)
	txt = strings.ReplaceAll(txt, "_", " ")
	if countWords(txt) < 2 {
		return ""
	}
	if strings.Contains(txt, "url") {
		txt = strings.Replace(txt, "url", "", -1)
		return fmt.Sprintf(" that indicates the endpoint URL for accessing to %s", txt)
	}
	return txt
}

func convertCamelCaseTo(str string) string {
	txt := strcase.SnakeCase(str)
	txt = strings.ReplaceAll(txt, "_", " ")
	if countWords(txt) < 2 {
		return ""
	}

	findPrefix := strings.ToLower(txt)

	if strings.HasPrefix(findPrefix, "get") {
		return replaceFirstWordGetWithRetrieve(txt)

	} else if strings.HasPrefix(findPrefix, "set") {
		return replaceFirstWordSetWithRetrieve(txt)

	} else if strings.HasPrefix(findPrefix, "init") {
		return replaceFirstWordSetWithInitialize(txt)

	} else if strings.HasPrefix(findPrefix, "delete") {

	} else if strings.HasPrefix(findPrefix, "is") {
		return replaceFirstWordSetWithCheckForIs(txt)

	} else if strings.HasPrefix(findPrefix, "create") {

	} else if strings.HasPrefix(findPrefix, "update") {

	} else if strings.HasPrefix(findPrefix, "has") {
		return replaceFirstWordSetWithCheckForHas(txt)
	} else if strings.HasPrefix(findPrefix, "handle") {

	} else if strings.HasPrefix(findPrefix, "process") {

	} else if strings.HasPrefix(findPrefix, "run") {

	} else if strings.HasPrefix(findPrefix, "load") {

	} else if strings.HasPrefix(findPrefix, "save") {

	} else if strings.HasPrefix(findPrefix, "init") {

	} else if strings.HasPrefix(findPrefix, "shutdown") {

	}

	return " which execute " + txt
}

func replaceFirstWordGetWithRetrieve(phrase string) string {
	return " that retrieve the" + phrase[3:]
}

func replaceFirstWordSetWithRetrieve(phrase string) string {
	return " which update the" + phrase[3:]
}

func replaceFirstWordSetWithInitialize(phrase string) string {
	return " to initializes the" + phrase[4:]
}

func replaceFirstWordSetWithCheckForIs(phrase string) string {
	return " to check the" + phrase[2:]
}

func replaceFirstWordSetWithCheckForHas(phrase string) string {
	return " to check the" + phrase[3:]
}

func countWords(sentence string) int {
	words := strings.Fields(sentence)
	return len(words)
}
