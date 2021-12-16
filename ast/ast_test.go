package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
	"testing"
)

func TestAST(t *testing.T) {
	fset := token.NewFileSet() // positions are relative to fset

	src := `package foo

import (
	"fmt"
)

func test(i int) (retI int, err error) {
	if i < 0 {
		return i, fmt.Errorf("negative number")
	} else {
		return i, nil
	}
}
`

	// Parse src but stop after processing the imports.
	f, err := parser.ParseFile(fset, "", src, parser.Trace|parser.ParseComments)
	fmt.Errorf("negative number")
	if err != nil {
		t.Error(err)
		return
	}

	funcs := make([]*ast.FuncDecl, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		return findFuctionsWithError(n, &funcs)
	})

	for _, fn := range funcs {
		ast.Inspect(fn, func(n ast.Node) bool {
			return findErrorReturnExpressions(fset, fn, n)
		})
	}
}

func findFuctionsWithError(n ast.Node, funcs *[]*ast.FuncDecl) bool {
	switch x := n.(type) {
	case *ast.FuncDecl:
		resultList := x.Type.Results.List
		if len(resultList) > 0 {
			lastResultType := resultList[len(resultList)-1].Type
			if ident, ok := lastResultType.(*ast.Ident); ok {
				if ident.Name == "error" {
					*funcs = append(*funcs, x)
				}
			}
		}
	}
	return true
}

func findErrorReturnExpressions(fset *token.FileSet, fnDecl *ast.FuncDecl, n ast.Node) bool {
	switch x := n.(type) {
	case *ast.ReturnStmt:
		resultList := x.Results
		if len(resultList) > 0 {
			lastExpr := resultList[len(resultList)-1]
			var buff strings.Builder
			//_ = ast.Fprint( os.Stdout, fset, lastExpr, nil)
			_ = printer.Fprint(&buff, fset, lastExpr)
			fmt.Printf("%s function returns error expression: %s\n", fnDecl.Name, buff.String())
		}
	}
	return true
}
