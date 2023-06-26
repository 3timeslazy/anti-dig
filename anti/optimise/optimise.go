// Copyright (c) 2023 Vladimir Fetisov
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

package optimise

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/ast/astutil"
)

func Optimise(fset *token.FileSet, file *ast.File, varTypes map[string]string) (*token.FileSet, *ast.File, error) {
	removeSelectors(file)
	renameVariables(file, fset, varTypes)
	removeImportsAliases(file)
	return fset, file, nil
}

func renameVariables(file *ast.File, fset *token.FileSet, varTypes map[string]string) {
	scope := map[string]bool{}
	replacements := map[string]string{}
	importMaxPos := map[string]token.Pos{}

	// During the first walk, we look for the last position
	// at which each imported package was used. Later on, this will
	// help us decide whether or not it's possible to shadow an import
	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		ident, ok := c.Node().(*ast.Ident)
		if !ok {
			return true
		}

		for _, imp := range file.Imports {
			if ident.Name == imp.Name.Name {
				importMaxPos[ident.Name] = ident.Pos()
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

		possibleNames := callExprNames(ident, file, varTypes)
		possibleNames = append(possibleNames, compositLitNames(ident)...)

		for i, name := range possibleNames {
			name = strcase.ToLowerCamel(name)

			// Make names with the length < 4 lower cased
			// because want to have name `db` instead of `dB` for the
			// type *DB
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

			// This if checks whether it's possible to shadow an import or
			// not. If the last line on which the import was used is the same
			// as current ident's line, that there will be the case:
			//
			//  a := a.New()
			//
			// In this case it's fine to shadow.
			// If the last line greater than ident's declaration than there will be
			// the case:
			//
			//  a := a.New()
			//  b := a.NewB() --> code won't compile because of the shadowing
			//
			if imppos, ok := importMaxPos[name]; ok {
				l1 := fset.Position(ident.Pos()).Line
				l2 := fset.Position(imppos).Line
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
		switch node := c.Node().(type) {
		case *ast.Ident:
			replaceIdent(node, replacements)

		case *ast.CompositeLit:
			for _, elt := range node.Elts {
				replaceIdent(elt, replacements)
			}

		case *ast.CallExpr:
			for _, arg := range node.Args {
				replaceIdent(arg, replacements)
			}
		}

		return true
	})
}

func compositLitNames(ident *ast.Ident) []string {
	assign, ok := ident.Obj.Decl.(*ast.AssignStmt)
	if !ok {
		return nil
	}
	rhs, ok := assign.Rhs[0].(*ast.CompositeLit)
	if !ok {
		return nil
	}

	names := []string{}

	switch rhs.Type.(type) {
	case *ast.SelectorExpr:
		sel := rhs.Type.(*ast.SelectorExpr)
		typ := sel.Sel.Name
		pkg := sel.X.(*ast.Ident).Name

		names = append(names, typ)
		names = append(names, pkg+"_"+typ)

	case *ast.ArrayType:
		arr := rhs.Type.(*ast.ArrayType)
		elem := arr.Elt.(*ast.SelectorExpr)
		elemType := elem.Sel.Name
		pkg := elem.X.(*ast.Ident).Name

		// Most likely an array is created for a specific dig group,
		// therefore the ident contains the group's name, which might be
		// a good variable name
		if i := strings.Index(ident.Name, "_"); i != -1 {
			names = append(names, ident.Name[i+1:])
		}

		names = append(names, "listOf_"+elemType)
		names = append(names, pkg+"_"+elemType+"_list")
	}

	return names
}

func callExprNames(ident *ast.Ident, file *ast.File, varTypes map[string]string) []string {
	assign, ok := ident.Obj.Decl.(*ast.AssignStmt)
	if !ok {
		return nil
	}
	call, ok := assign.Rhs[0].(*ast.CallExpr)
	if !ok {
		return nil
	}

	varType := varTypes[ident.Name]
	if ident.Name == "err" {
		return nil
	}

	fn := call.Fun.(*ast.SelectorExpr)
	pkg := ""

	// Function calls generated by anti-dig usually have two different forms:
	//  1. Provider as a named function: package.Provider(args)
	//  2. Provider as an anonymous function: package.Provider.func1(args)
	switch fn.X.(type) {
	case *ast.SelectorExpr:
		anonFn := fn.X.(*ast.SelectorExpr)
		pkg = anonFn.X.(*ast.Ident).Name
		fn = anonFn

		cg := &ast.CommentGroup{
			List: []*ast.Comment{{
				Text:  "// TODO: anonymous functions are not yet supported. Please, fix it manually",
				Slash: ident.Pos() - 1,
			}},
		}

		file.Comments = append(file.Comments, cg)

	case *ast.Ident:
		pkg = fn.X.(*ast.Ident).Name
	}

	fnName := fn.Sel.Name

	fnName, _ = strings.CutSuffix(fnName, "Provider")
	fnName, _ = strings.CutPrefix(fnName, "New")

	names := []string{varType}
	names = append(names, pkg+"_"+varType)
	names = append(names, fnName)
	names = append(names, pkg+"_"+fnName)

	return names
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

	// strings.ReplaceAll instead of raw assignment helps in
	// a situation when the ident isn't just a variable, but variable with
	// a selector, like `var1.Field`
	ident.Name = strings.Replace(ident.Name, name, repl, 1)
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
		assign, ok := c.Node().(*ast.AssignStmt)
		if !ok {
			return true
		}
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
