// Copyright (c) 2021 Vladimir Fetisov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package optimiser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/ast/astutil"
)

type Optimiser struct {
	Output   io.Writer
	VarTypes map[string]string
}

func (opt *Optimiser) PrintOptimised(generated string) error {
	fset := token.NewFileSet()

	generatedFile, err := parser.ParseFile(fset, "provide.gen.go", generated, parser.ParseComments)
	if err != nil {
		return err
	}

	removeSelectors(generatedFile)
	renameVariables(generatedFile, fset, opt.VarTypes)
	removeImportsAliases(generatedFile)

	err = printer.Fprint(opt.Output, fset, generatedFile)
	if err != nil {
		return err
	}

	return err
}

func renameVariables(file *ast.File, fset *token.FileSet, varTypes map[string]string) {
	scope := map[string]bool{}
	replacements := map[string]string{}
	imports := map[string]token.Pos{}

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		node := c.Node()

		switch ident := node.(type) {
		case *ast.Ident:
			if slices.ContainsFunc(file.Imports, func(imp *ast.ImportSpec) bool {
				return ident.Name == imp.Name.Name
			}) {
				imports[ident.Name] = ident.Pos()
			}
		}

		return true
	})

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		ident, ok := c.Node().(*ast.Ident)
		if !ok {
			return true
		}

		// If ident.Name has been seen in the scope before
		// it means that we've already found a replacement for this
		// therefore it makes no sense to proceed
		if scope[ident.Name] {
			return true
		}
		if ident.Obj == nil || ident.Obj.Kind != ast.Var {
			return true
		}

		// TODO: later
		vtype := varTypes[ident.Name]
		if vtype == "" {
			return true
		}

		possibleNames := []string{vtype}

		assign, ok := ident.Obj.Decl.(*ast.AssignStmt)
		if ok {
			call := assign.Rhs[0].(*ast.CallExpr)
			fun := call.Fun.(*ast.SelectorExpr)
			pkg := fun.X.(*ast.Ident)

			fnName := fun.Sel.Name
			fnName, _ = strings.CutSuffix(fnName, "Provider")
			fnName, _ = strings.CutPrefix(fnName, "New")

			possibleNames = append(possibleNames, pkg.Name+"_"+vtype)  // pkgType
			possibleNames = append(possibleNames, fnName)              // funcName
			possibleNames = append(possibleNames, pkg.Name+"_"+fnName) // pkgFuncName
		}

		for i, name := range possibleNames {
			name = strcase.ToLowerCamel(name)
			if len(name) < 4 {
				name = strings.ToLower(name)
			}

			possibleNames[i] = name
		}
		slices.SortFunc(possibleNames, func(a, b string) bool {
			return len(a) < len(b)
		})

		for _, name := range possibleNames {
			if token.IsKeyword(name) || !token.IsIdentifier(name) {
				continue
			}
			if scope[name] {
				continue
			}
			if impPos, ok := imports[name]; ok {
				l1 := fset.Position(ident.Pos()).Line
				l2 := fset.Position(impPos).Line
				if l2 > l1 {
					continue
				}
			}

			scope[name] = true
			replacements[ident.Name] = name
			break
		}

		scope[ident.Name] = true
		return true
	})

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		switch ident := c.Node().(type) {
		case *ast.Ident:
			if ident.Obj == nil || ident.Obj.Kind != ast.Var {
				break
			}

			replaceIdent(ident, replacements)
		}

		switch comp := c.Node().(type) {
		case *ast.CompositeLit:
			for _, elt := range comp.Elts {
				replaceIdent(elt, replacements)
			}
		}

		switch call := c.Node().(type) {
		case *ast.CallExpr:
			for _, arg := range call.Args {
				replaceIdent(arg, replacements)
			}
		}

		return true
	})
}

func replaceIdent(node ast.Node, replacements map[string]string) {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return
	}

	name := strings.Split(ident.Name, ".")[0]
	repl, ok := replacements[name]
	if !ok {
		return
	}

	ident.Name = strings.ReplaceAll(ident.Name, name, repl)
}

func removeImportsAliases(file *ast.File) {
	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		pkg, ok := c.Node().(*ast.ImportSpec)
		if !ok {
			return true
		}
		if strings.HasSuffix(pkg.Path.Value, "/"+pkg.Name.Name+`"`) {
			pkg.Name = nil
		}

		return true
	})
}

func removeSelectors(file *ast.File) {
	replacements := map[string]string{}

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		switch assign := c.Node().(type) {
		case *ast.AssignStmt:
			if len(assign.Rhs) > 1 {
				panic("functions with more than one return argument aren't supported yet")
			}

			varname := assign.Lhs[0].(*ast.Ident).Name

			switch selector := assign.Rhs[0].(type) {
			case *ast.SelectorExpr:
				repl := fmt.Sprintf("%s.%s", selector.X, selector.Sel)
				replacements[varname] = repl

				c.Delete()
			}
		}

		return true
	})

	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		switch ident := c.Node().(type) {
		case *ast.Ident:
			if ident.Obj == nil || ident.Obj.Kind != ast.Var {
				break
			}

			repl, ok := replacements[ident.Name]
			if !ok {
				break
			}

			ident.Name = repl
		}

		return true
	})
}
