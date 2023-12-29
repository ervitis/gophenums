package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

// using ast
func main() {
	paths := make(map[string]struct{})
	absPath, _ := filepath.Abs(".")
	if err := filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			if _, ok := paths[path]; ok {
				return nil
			}
			paths[path] = struct{}{}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	fset := token.NewFileSet()

	typeNames := make([]string, 0)

	for fi := range paths {
		f, err := parser.ParseFile(fset, fi, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		ast.Inspect(f, func(node ast.Node) bool {
			if node == nil {
				return false
			}
			switch x := node.(type) {
			case *ast.GenDecl:
				if len(x.Specs) == 0 || x.Doc == nil {
					return true
				}
				for _, commentGroup := range x.Doc.List {
					if !strings.Contains(commentGroup.Text, "gophenum:generate") {
						return true
					}
					if tSpec, ok := x.Specs[0].(*ast.TypeSpec); ok {
						typeNames = append(typeNames, tSpec.Name.Name)
					}
				}
			}
			return true
		})
	}
}
