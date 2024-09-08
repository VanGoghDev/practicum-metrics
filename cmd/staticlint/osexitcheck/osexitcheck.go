// Анализатор, который проверяет, что в функции main не вызывается os.exit.
package osexitcheck

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer - instance for multichecker.
var Analyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.exit in main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					return true
				}
				return false
			case *ast.ExprStmt:
				if call, ok := x.X.(*ast.CallExpr); ok {
					if isFuncFromPkg(call.Fun, "os", "Exit") {
						pass.Reportf(x.Pos(), "preventing call os.Exit in main function")
					}
				}
				return true
			default:
				return true
			}
		})
	}
	//nolint:all // using code style as described in analyzer documentation
	return nil, nil
}

func isFuncFromPkg(expr ast.Expr, pkg, name string) bool {
	s, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(s.X, pkg) && isIdent(s.Sel, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}
