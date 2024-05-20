package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"moul.io/http2curl"
	"net/http"
	"strings"

	"github.com/fatih/astrewrite"
	"github.com/stoewer/go-strcase"
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
	if file.cfg.Signature != nil {
		return "//\n// Author: " + *file.cfg.Signature + "."
	}
	return ""
}

func (file *file) autoComment() ([]byte, error) {
	var comments []*ast.CommentGroup

	for _, decl := range file.f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {

			switch genDecl.Tok {
			case token.TYPE:
				file.commentType(genDecl)
			case token.CONST:
				file.commentConst(genDecl)
			case token.VAR:
				file.commentVar(genDecl)
			default:
			}
		}

		if genDecl, ok := decl.(*ast.FuncDecl); ok {
			if genDecl.Doc.Text() == "" && genDecl.Name.Name != "main" && genDecl.Name.Name != "init" {
				txt := file.commentFunc(genDecl)
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

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (file *file) callOpenAI(functionCode string) (string, error) {
	prompt := fmt.Sprintf("Generate a detailed comment in English for the following Go function. The comment should be written in a way that is helpful for other developers. Include the purpose of the function, a description of its parameters and return values, potential error conditions, and any side effects or important details. Here is the function :\n%s", functionCode)

	requestBody, err := json.Marshal(map[string]interface{}{
		"max_tokens":  150,
		"temperature": 0.7,
		"model":       "gpt-3.5-turbo",
		"messages": []OpenAIMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: prompt},
		},
	})

	if err != nil {
		return "", fmt.Errorf("error creating request body: %v", err)
	}

	req, err := http.NewRequest("POST", file.cfg.OpenAIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+*file.cfg.OpenAIAPIKey)

	client := &http.Client{}

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(fmt.Sprintf("%s", command))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if text, ok := choice["text"].(string); ok {
				return text, nil

			} else {
				return "", errors.New("error: no text found in response choice")
			}
		} else {
			return "", errors.New("error: invalid choice format")
		}
	}

	return "", errors.New("error: no choices found in response")
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
					errorReturnMsg = "// It's return an error if the initialization fails, otherwise nil.\n"
				}
			}
		}
	}

	txt := fmt.Sprintf("// %s creates a new instance%s.\n", fn.Name.Name, instanceReturnMsg)
	if (fn.Type.Params == nil || len(fn.Type.Params.List) == 0) && (fn.Type.Results == nil || len(fn.Type.Results.List) == 0) {
	} else {
		if fn.Type.Params != nil {
			if len(fn.Type.Params.List) != 0 {
				txt += fmt.Sprintf("// It initializes the %s with the provided ", initializesMsg)
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
		txt += ".\n"
	}
	if errorReturnMsg != "" {
		txt += errorReturnMsg
	}

	return txt
}

func (file *file) isOpenAIActive() bool {
	if file.cfg.OpenAIActive == nil || *file.cfg.OpenAIActive == false {
		return false
	}
	if file.cfg.OpenAIAPIKey == nil || *file.cfg.OpenAIAPIKey == "" {
		log.Fatal("Please set your OpenAI API key in the OPENAI_API_KEY variable.")
	}
	if file.cfg.OpenAIURL == "" {
		log.Fatal("Please set the OpenAI API URL in the OPENAI_API_URL variable.")
	}
	return true
}

func (file *file) commentConstWithOpenAI(genDecl *ast.GenDecl) {

}

func (file *file) commentConst(genDecl *ast.GenDecl) {
	if file.isOpenAIActive() {
		file.commentConstWithOpenAI(genDecl)
		return
	}
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

					explainConst := convertVarToCamelCaseTo(name.Name)
					txt := fmt.Sprintf("// %s is a %sconstant%s.", name.Name, exported, explainConst)

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
}

func (file *file) commentVarWithOpenAI(genDecl *ast.GenDecl) {

}

func (file *file) commentVar(genDecl *ast.GenDecl) {
	if file.isOpenAIActive() {
		file.commentVarWithOpenAI(genDecl)
		return
	}
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

				explainVar := convertVarToCamelCaseTo(name.Name)

				if decl, ok := name.Obj.Decl.(*ast.ValueSpec); ok {
					txt := fmt.Sprintf("// %s is a %svariable of type %s%s.", name.Name, exported, fmt.Sprintf("%s", decl.Type), explainVar)

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
}

func (file *file) commentTypeWithOpenAI(genDecl *ast.GenDecl) {

}

func (file *file) commentType(genDecl *ast.GenDecl) {
	if file.isOpenAIActive() {
		file.commentTypeWithOpenAI(genDecl)
		return
	}
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
}

func isNewFunc(name string) bool {
	return strings.HasPrefix(name, "New")
}

func (file *file) commentFuncWithOpenAI(fn *ast.FuncDecl) string {
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

	functionCode := fmt.Sprintf(`func %s%s(%s) (%s) {
    // Function implementation
	}`, funcType, functionName, input, output)

	txt, err := file.callOpenAI(functionCode)
	if err != nil {
		fmt.Println(fmt.Sprintf("get error on openAI %v", err))
		return ""
	}
	return txt
}

func (file *file) commentFunc(fn *ast.FuncDecl) string {
	if file.isOpenAIActive() {
		return file.commentFuncWithOpenAI(fn)
	}

	var (
		txt     string
		inputs  []ast.Expr
		outputs []ast.Expr
	)
	if isNewFunc(fn.Name.Name) {
		return newFuncTxt(fn)
	}

	privateValue := ""
	if !fn.Name.IsExported() {
		privateValue = "private "
	}

	explainFunc := convertCamelCaseTo(fn.Name.Name)

	if (fn.Type.Params == nil || len(fn.Type.Params.List) == 0) && (fn.Type.Results == nil || len(fn.Type.Results.List) == 0) {
		if fn.Recv != nil {
			funcType := ""
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				if t, ok := fn.Recv.List[0].Type.(*ast.Ident); ok {
					funcType = t.String()
				}
			}
			txt = fmt.Sprintf("// %s is a %smethod%s that belongs to the %s struct.\n// It does not take any arguments.\n", fn.Name.Name, privateValue, explainFunc, funcType)
		} else {
			txt = fmt.Sprintf("// %s is a %smethod%s.\n// It does not take any arguments.\n", fn.Name.Name, privateValue, explainFunc)
		}

	} else {
		if fn.Recv != nil {
			funcType := ""
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				if t, ok := fn.Recv.List[0].Type.(*ast.Ident); ok {
					funcType = t.String()
				}
			}
			txt = fmt.Sprintf("// %s is a %smethod%s that belongs to the %s struct", fn.Name.Name, privateValue, explainFunc, funcType)
		} else {
			txt = fmt.Sprintf("// %s is a %smethod%s", fn.Name.Name, privateValue, explainFunc)
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
					inputs = append(inputs, param.Type)
					txt += fmt.Sprintf("%s %s of type %s", indefiniteArticle(fmt.Sprintf("%s", name.Name)), name.Name, getTypeName(param.Type))
				}
			}
		}

		if fn.Type.Results != nil && len(fn.Type.Results.List) != 0 {
			errorReturnMsg := ""
			for i, res := range fn.Type.Results.List {
				typeReturnKey := fmt.Sprintf("%s", res.Type)
				if typeReturnKey == "error" {
					errorReturnMsg = ".\n// It's return an error if fails, otherwise nil"
					outputs = append(outputs, res.Type)
				} else {
					if i > 0 {
						txt += " and "
					} else {
						txt += "\n// and returns "
					}
					outputs = append(outputs, res.Type)
					txt += fmt.Sprintf("%s %s", indefiniteArticle(fmt.Sprintf("%s", res.Type)), res.Type)
				}
			}
			txt += errorReturnMsg
		}
		txt += ".\n"
	}

	txt += exampleGenerator(fn.Name.Name, inputs, outputs)

	return txt
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

func exampleGenerator(funcName string, inputs []ast.Expr, outputs []ast.Expr) string {
	if len(inputs) == 0 && len(outputs) == 0 {
		return ""
	}

	exampleComment := "//\n// Example:\n//   "

	hasError := false

	if len(inputs) > 0 && detectExprTypeKey(inputs[0]) == "ctx" {
		exampleComment += "ctx := context.Background()\n//   "
	}

	outputsStr := make([]string, len(outputs))
	var outputsStrWithoutErrors []string
	// Générer des exemples pour les paramètres de sortie
	if len(outputs) > 0 {
		for i, output := range outputs {
			exprType := detectExprTypeKey(output)
			if exprType == "unknown" {
				return ""
			} else if exprType == "err" {
				hasError = true
			} else {
				outputsStrWithoutErrors = append(outputsStrWithoutErrors, fmt.Sprintf("%v", exprType))
			}
			outputsStr[i] = fmt.Sprintf("%v", exprType)
		}
		exampleComment += strings.Join(outputsStr, ", ")
		exampleComment += " := "
	}

	exampleComment += funcName + "("

	// Générer des exemples pour les paramètres d'entrée
	inputsStr := make([]string, len(inputs))
	for i, input := range inputs {
		exprType := detectExprTypeValue(input)
		if exprType == "unknown" {
			return ""
		}
		inputsStr[i] = fmt.Sprintf("%s", exprType)
	}
	exampleComment += strings.Join(inputsStr, ", ")
	exampleComment += ")\n"

	if hasError {
		exampleComment += "//   if err != nil {\n//       log.Fatalf(\"Error: %v\", err)\n//   }\n"
	}

	lenOutputs := len(outputsStrWithoutErrors)
	if lenOutputs > 0 {
		exampleComment += "//   fmt.Printf(\"" + generatePrintfFormat(lenOutputs) + "\", " + strings.Join(outputsStrWithoutErrors, ", ") + ")\n"
	}
	return exampleComment
}
func generatePrintfFormat(sliceLength int) string {
	if sliceLength <= 0 {
		return ""
	}
	return strings.Repeat("%v ", sliceLength)
}

func detectExprTypeKey(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		switch v.Name {
		case "bool":
			return "valid"
		case "int":
			return "nb"
		case "string":
			return "str"
		case "float32":
			return "nb"
		case "float64":
			return "nb"
		case "error":
			return "err"
		case "Context":
			return "ctx"
		default:
			fmt.Println(fmt.Sprintf("detectExprTypeKey : %+v", v))
			return "unknown"
		}
	case *ast.SelectorExpr:
		// handle qualified types like "pkg.Type"
		return detectExprTypeValue(v.Sel)
	default:
		return "unknown"
	}
}

func detectExprTypeValue(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		switch v.Name {
		case "bool":
			return "true"
		case "int":
			return "50"
		case "string":
			return "my-string"
		case "float32":
			return "56.32"
		case "float64":
			return "56.64"
		case "error":
			return "nil"
		case "Context":
			return "ctx"
		default:
			fmt.Println(fmt.Sprintf("%+v", v))
			return "unknown"
		}
	case *ast.SelectorExpr:
		// handle qualified types like "pkg.Type"
		return detectExprTypeValue(v.Sel)
	default:
		return "unknown"
	}
}
