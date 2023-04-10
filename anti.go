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

package antidig

import (
	"container/list"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type AntiDig struct {
	output io.Writer
	exprs  []string

	callstack    *list.List
	fnsArgs      map[string][]string
	fnsVars      map[string][]string
	fnsSuffixies map[string][]string

	flattenVars map[string]bool

	typeVarname    map[reflect.Type]string
	typeVarnameSeq int
	typeSeqname    map[reflect.Type]int
	pkgAlias       map[string]string
}

var Anti = AntiDig{
	output: os.Stdout,

	callstack:    list.New(),
	fnsArgs:      map[string][]string{},
	fnsVars:      map[string][]string{},
	fnsSuffixies: map[string][]string{},

	flattenVars: map[string]bool{},

	typeVarname:    map[reflect.Type]string{},
	typeVarnameSeq: 0,
	typeSeqname:    map[reflect.Type]int{},
	pkgAlias:       map[string]string{},
}

func (anti *AntiDig) Generate() {
	decls := []string{
		"package main\n",
		anti.generateImports(),
		anti.generateFunc(),
	}

	for _, decl := range decls {
		fmt.Fprintln(anti.output, decl)
	}
}

func (anti *AntiDig) generateFunc() string {
	out := "func main() {\n"

	for _, expr := range anti.exprs {
		out += fmt.Sprintf("\t%s\n", expr)
	}

	out += "}\n"

	return out
}

func (anti *AntiDig) generateImports() string {
	out := "import (\n"

	for pkg, alias := range anti.pkgAlias {
		out += fmt.Sprintf("\t%s \"%s\"\n", alias, pkg)
	}

	out += ")\n"

	return out
}

func (anti *AntiDig) PushFnCall(fnName string) {
	anti.callstack.PushBack(fnName)
}

func (anti *AntiDig) PopFnCall() {
	elem := anti.callstack.Back()
	anti.callstack.Remove(elem)
}

func (anti *AntiDig) AppendFnArg(arg string) {
	curr := anti.currFn()
	anti.fnsArgs[curr] = append(anti.fnsArgs[curr], arg)
}

func (anti *AntiDig) AppendFnVar(varname string) {
	curr := anti.currFn()
	anti.fnsVars[curr] = append(anti.fnsVars[curr], varname)
}

func (anti *AntiDig) AppendFnSuffix(suffix string) {
	curr := anti.currFn()
	anti.fnsSuffixies[curr] = append(anti.fnsSuffixies[curr], suffix)
}

func (anti *AntiDig) FnSuffixes() []string {
	return anti.fnsSuffixies[anti.currFn()]
}

func (anti *AntiDig) FnArgs() []string {
	return anti.fnsArgs[anti.currFn()]
}

func (anti *AntiDig) FnVars() []string {
	return anti.fnsVars[anti.currFn()]
}

func (anti *AntiDig) AddFlatten(varname string, flatten bool) {
	anti.flattenVars[varname] = flatten
}

func (anti *AntiDig) Flatten(varname string) bool {
	return anti.flattenVars[varname]
}

func (anti *AntiDig) TypeVarname(typ reflect.Type) string {
	varname, ok := anti.typeVarname[typ]
	if ok {
		return varname
	}

	anti.typeVarnameSeq++
	varname = fmt.Sprintf("var%d", anti.typeVarnameSeq)
	anti.typeVarname[typ] = varname

	return varname
}

func (anti *AntiDig) TypeSeqname(typ reflect.Type) string {
	varname := anti.TypeVarname(typ)

	_, ok := anti.typeSeqname[typ]
	if !ok {
		anti.typeSeqname[typ] = 0
		return fmt.Sprintf("%s_0", varname)
	}

	anti.typeSeqname[typ]++
	return fmt.Sprintf("%s_%d", varname, anti.typeSeqname[typ])
}

func (anti *AntiDig) currFn() string {
	if anti.callstack.Len() == 0 {
		return ""
	}
	return anti.callstack.Back().Value.(string)
}

func (anti *AntiDig) Print(expr ...string) {
	anti.exprs = append(anti.exprs, expr...)
}

func (anti *AntiDig) PkgAlias(pkgname string) string {
	alias, ok := anti.pkgAlias[pkgname]
	if ok {
		return alias
	}

	path := strings.Split(pkgname, "/")
	alias = path[len(path)-1]
	anti.pkgAlias[pkgname] = alias
	return alias
}

var panicExpr = []string{"if err != nil {", "\tpanic(err)", "}"}
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
