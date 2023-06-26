package overwrite

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"unicode"

	rename "github.com/3timeslazy/anti-dig/anti/gorename"

	"golang.org/x/tools/go/ast/astutil"
)

func Rename(fset *token.FileSet, file *ast.File) error {
	pkgAliases := map[string]string{}

	imports := astutil.Imports(fset, file)
	for _, imp := range imports {
		for _, spec := range imp {
			pkgAliases[spec.Name.Name] = spec.Path.Value
		}
	}

	froms, tos := []string{}, []string{}

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		node := c.Node()

		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkg := sel.X.(*ast.Ident).Name
		fnName := sel.Sel.Name
		if unicode.IsLower(rune(fnName[0])) {
			newFnName := fmt.Sprintf("%c%s", unicode.ToUpper(rune(fnName[0])), fnName[1:])

			froms = append(froms, fmt.Sprintf("%s.%s", pkgAliases[pkg], fnName))
			tos = append(tos, newFnName)
			sel.Sel.Name = newFnName
		}

		return true
	})

	return rename.Main(&build.Default, froms, tos)
}
