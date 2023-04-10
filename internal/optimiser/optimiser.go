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
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

type Optimiser struct {
	selRepls map[string]string
}

func New() *Optimiser {
	return &Optimiser{
		selRepls: map[string]string{},
	}
}

func (opt *Optimiser) PrintOptimised(src string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "main.gen.go", src, 0)
	if err != nil {
		return err
	}

	// Gather information
	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		node := c.Node()

		switch node := node.(type) {
		case *ast.AssignStmt:
			if len(node.Rhs) > 1 {
				panic("len(Rhs) > 1 currently unsupported")
			}

			varname := node.Lhs[0].(*ast.Ident).Name

			switch rhs := node.Rhs[0].(type) {
			case *ast.SelectorExpr:
				repl := fmt.Sprintf("%s.%s", rhs.X, rhs.Sel)
				opt.selRepls[varname] = repl

				c.Delete()

			default:
			}

		default:
		}

		return true
	})

	// Replace variables
	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
		node := c.Node()

		switch ident := node.(type) {
		case *ast.Ident:
			if ident.Obj == nil || ident.Obj.Kind != ast.Var {
				break
			}

			repl, ok := opt.selRepls[ident.Name]
			if !ok {
				break
			}

			ident.Name = repl

			c.Replace(ident)
		}

		return true
	})

	err = printer.Fprint(os.Stdout, fset, file)
	if err != nil {
		return err
	}

	return err
}
