package osexit

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzer prohibits using a direct os.Exit call in the main function of the main package
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  `check for os.Exit in main func`,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if f, ok := node.(*ast.FuncDecl); ok && f.Name.Name == "main" {
				checkOsExit(f, pass)
			}
			return true
		})
	}
	return nil, nil
}

func checkOsExit(f *ast.FuncDecl, pass *analysis.Pass) {
	ast.Inspect(f, func(node ast.Node) bool {
		if c, ok := node.(*ast.CallExpr); ok {
			if s, ok := c.Fun.(*ast.SelectorExpr); ok {
				if strings.Contains("os.Exit", s.Sel.Name) {
					pass.Reportf(c.Pos(), "os.Exit found in main function")
				}
			}
		}
		return true
	})
}
