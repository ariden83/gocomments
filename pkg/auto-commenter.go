package pkg

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/astrewrite"
)

// AutoCommentDir ...
func AutoCommentDir(dir string) {
	pkg, err := build.ImportDir(dir, 0)
	autoCommentImportedPkg(pkg, err)
}

func autoCommentImportedPkg(pkg *build.Package, err error) {
	log.Println("autoCommentImportedPkg")
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			return
		}
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	files := make([]string, 0)

	files = append(files, pkg.GoFiles...)
	if pkg.Dir != "." {
		for i, f := range files {
			files[i] = filepath.Join(pkg.Dir, f)
		}
	}

	readingFiles(files...)
}

func AutoCommentFiles(files ...string) {
	readingFiles(files...)
}

func readingFiles(files ...string) {
	log.Println(fmt.Sprintf("readingFiles %+v", files))
	fileBodyMap := make(map[string][]byte)

	for _, file := range files {
		fileBody, err := ioutil.ReadFile(file)
		if err != nil {
			log.Println("+++++++++", err)
			continue
		}

		fileBodyMap[file] = fileBody
	}

	autoCmntr := AutoCommenter{}
	_ = autoCmntr.AutoCommentFiles(fileBodyMap)
}

// AutoCommenter ...
type AutoCommenter struct{}

type pkg struct {
	fileSet   *token.FileSet
	files     map[string]*file
	typesPkg  *types.Package
	typesInfo *types.Info
}

type file struct {
	pkg      *pkg
	f        *ast.File
	fset     *token.FileSet
	src      []byte
	filename string
}

// AutoCommentFiles ...
func (auto *AutoCommenter) AutoCommentFiles(filesMap map[string][]byte) error {
	pkg := &pkg{
		fileSet: token.NewFileSet(),
		files:   make(map[string]*file),
	}

	var packageName string

	for fileName, body := range filesMap {
		f, err := parser.ParseFile(pkg.fileSet, fileName, body, parser.ParseComments)
		if err != nil {
			return err
		}

		if packageName == "" {
			packageName = f.Name.Name
		} else if f.Name.Name != packageName {
			return fmt.Errorf("%s is in package %s, not %s", fileName, f.Name.Name, packageName)
		}

		pkg.files[fileName] = &file{
			pkg:      pkg,
			f:        f,
			fset:     pkg.fileSet,
			src:      body,
			filename: fileName,
		}
	}

	if len(pkg.files) != 0 {
		return pkg.autoComment()
	}

	return nil
}

func (pkg *pkg) autoComment() error {
	for _, file := range pkg.files {
		file.autoComment()
	}
	return nil
}

func (file *file) autoComment() {
	if strings.HasSuffix(file.filename, "_test.go") {
		return
	}

	var comments []*ast.CommentGroup

	ast.Inspect(file.f, func(node ast.Node) bool {
		if node == nil {
			return false
		}

		c, ok := node.(*ast.CommentGroup)
		if ok {
			comments = append(comments, c)
		}

		switch fn := node.(type) {
		case *ast.ValueSpec:
			for _, name := range fn.Names {
				if name.IsExported() && fn.Doc.Text() == "" {
					if decl, ok := name.Obj.Decl.(*ast.ValueSpec); ok {
						txt := fmt.Sprintf("// %s is a variable of type %s which provides .", name.Name, fmt.Sprintf("%s", decl.Type))

						fn.Doc = &ast.CommentGroup{
							List: []*ast.Comment{{
								Text:  txt,
								Slash: fn.Pos() - 1,
							}},
						}
					}
				}
			}
		case *ast.StructType:
			log.Printf("// Structure: %+v\n", fn.Incomplete)
			log.Printf("// Structure: %+v\n", fn.Struct)

			for _, field := range fn.Fields.List {

				txt := "// test"
				for _, name := range field.Names {
					txt += fmt.Sprintf("// - %s\n", name.Name)
				}

				field.Doc = &ast.CommentGroup{
					List: []*ast.Comment{{
						Text:  txt,
						Slash: fn.Pos() - 1,
					}},
				}
			}

		case *ast.FuncDecl:
			if fn.Name.IsExported() && fn.Doc.Text() == "" {
				fn.Doc = &ast.CommentGroup{
					List: []*ast.Comment{{
						Text:  getFuncComments(fn),
						Slash: fn.Pos() - 1,
					}},
				}
			}
		}

		return true
	})

	file.f.Comments = comments

	reWriteFunc := func(node ast.Node) (ast.Node, bool) {

		return node, true
	}

	newAst := astrewrite.Walk(file.f, reWriteFunc)
	var buf bytes.Buffer
	_ = printer.Fprint(&buf, file.fset, newAst)

	f, err := os.OpenFile(file.filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.Seek(0, 0)

	f.Write(buf.Bytes())
	f.Sync()

}

func getFuncComments(fn *ast.FuncDecl) string {
	txt := ""
	if (fn.Type.Params == nil || len(fn.Type.Params.List) == 0) && (fn.Type.Results == nil || len(fn.Type.Results.List) == 0) {
		if fn.Recv != nil {
			txt = fmt.Sprintf("// %s is a method that belongs to the %s struct.\n// It does not take any arguments.\n", fn.Name.Name, fn.Recv.List[0].Type.(*ast.Ident))
		} else {
			txt = fmt.Sprintf("// %s is a method .\n// It does not take any arguments.\n", fn.Name.Name)
		}

	} else {
		if fn.Recv != nil {
			txt = fmt.Sprintf("// %s is a method that belongs to the %s struct", fn.Name.Name, fn.Recv.List[0].Type.(*ast.Ident))
		} else {
			txt = fmt.Sprintf("// %s is a method", fn.Name.Name)
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

		if fn.Type.Results != nil {
			if len(fn.Type.Results.List) != 0 {
				txt += "\n// and returns "
			}
			for i, res := range fn.Type.Results.List {
				if i > 0 {
					txt += " and "
				}
				txt += fmt.Sprintf("%s %s", IndefiniteArticle(fmt.Sprintf("%s", res.Type)), res.Type)
			}
			txt += "."
		}
		txt += "\n"
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

type functionSpec struct {
	Name   string
	Prefix string
	Kind   string
}

func (file *file) isLintedFuncDoc(fn *ast.FuncDecl) (*functionSpec, error) {
	if !ast.IsExported(fn.Name.Name) {

		return nil, nil
	}
	kind := "function"
	name := fn.Name.Name
	prefix := fn.Name.Name + " "
	if fn.Doc == nil {
		return &functionSpec{
			Name:   name,
			Prefix: prefix,
			Kind:   kind,
		}, fmt.Errorf("exported %s %s should have comment or be unexported", kind, name)
	}
	s := fn.Doc.Text()

	if !strings.HasPrefix(s, prefix) {
		return &functionSpec{
			Name:   name,
			Prefix: prefix,
			Kind:   kind,
		}, fmt.Errorf(`comment on exported %s %s should be of the form "%s..."`, kind, name, prefix)
	}
	return nil, nil
}
