package main

import (
	"fmt"
	"go/ast"
	"strings"
)

type defaultProcess struct {
	activeExamples bool
}

func (d *defaultProcess) isActive() bool {
	return true
}

func (d *defaultProcess) commentFunc(fn *ast.FuncDecl) (string, error) {

	var (
		txt     string
		inputs  []ast.Expr
		outputs []ast.Expr
	)
	if isNewFunc(fn.Name.Name) {
		return d.newFuncTxt(fn), nil
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

	txt += d.exampleGenerator(fn.Name.Name, inputs, outputs)

	return txt, nil
}

func (d *defaultProcess) commentConst(name string, exported bool) (string, error) {
	var exportedTxt string
	if !exported {
		exportedTxt = "private "
	}
	explainConst := convertVarToCamelCaseTo(name)
	txt := fmt.Sprintf("// %s is a %sconstant%s.", name, exportedTxt, explainConst)

	return txt, nil
}

func (d *defaultProcess) commentVar(name, declType, explainVar string, exported bool) (string, error) {
	var exportedTxt string
	if !exported {
		exportedTxt = "private "
	}
	txt := fmt.Sprintf("// %s is a %svariable of type %s%s.", name, exportedTxt, declType, explainVar)

	return txt, nil
}

func (d *defaultProcess) commentType(genDecl *ast.GenDecl) (string, error) {
	return "", nil
}

func (d *defaultProcess) newFuncTxt(fn *ast.FuncDecl) string {
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

func (d *defaultProcess) exampleGenerator(funcName string, inputs []ast.Expr, outputs []ast.Expr) string {
	if !d.activeExamples {
		return ""
	} else if len(inputs) == 0 && len(outputs) == 0 {
		return ""
	}

	exampleComment := "//\n// Example:\n//   "

	hasError := false

	if len(inputs) > 0 && detectExprTypeKey(inputs[0]) == "ctx" {
		exampleComment += "ctx := context.Background()\n//   "
	}

	outputsStr := make([]string, len(outputs))
	var outputsStrWithoutErrors []string
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
