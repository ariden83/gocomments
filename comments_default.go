package main

import (
	"fmt"
	"go/ast"
)

type defaultProcess struct{}

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

	txt += exampleGenerator(fn.Name.Name, inputs, outputs)

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
