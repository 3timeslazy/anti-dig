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

	"github.com/3timeslazy/anti-dig/internal/optimiser"

	"github.com/iancoleman/strcase"
)

type AntiDig struct {
	output io.Writer
	exprs  []string

	callstack    *list.List
	fnsArgs      map[string][]string
	fnsVars      map[string][]string
	fnsSuffixies map[string][]string

	flattenVars map[string]bool

	varnames      map[typeAlias]string
	groupVarnames map[typeAlias]string
	varnameSeq    int
	seqnames      map[typeAlias]int

	pkgAlias      map[string]string
	allPkgAliases map[string]bool

	optimiser *optimiser.Optimiser
}

type typeAlias struct {
	Alias string
	Type  reflect.Type
}

func NewAnti(output io.Writer) *AntiDig {
	return &AntiDig{
		callstack:    list.New(),
		fnsArgs:      map[string][]string{},
		fnsVars:      map[string][]string{},
		fnsSuffixies: map[string][]string{},

		flattenVars: map[string]bool{},

		varnames:      map[typeAlias]string{},
		groupVarnames: map[typeAlias]string{},
		varnameSeq:    0,
		seqnames:      map[typeAlias]int{},

		pkgAlias:      map[string]string{},
		allPkgAliases: map[string]bool{},

		optimiser: optimiser.New(output),
	}
}

var Anti = NewAnti(os.Stdout)

func (anti *AntiDig) Generate(invokedType reflect.Type) error {
	decls := []string{
		"package main\n",
		anti.generateImports(),
		anti.generateFunc(invokedType),
	}

	return anti.optimiser.PrintOptimised(strings.Join(decls, "\n"))
}

func (anti *AntiDig) generateFunc(invokedType reflect.Type) string {
	anti.exprs = append(anti.exprs, anti.returnStmt(invokedType))

	returnedTypes := anti.returnedTypes(invokedType)
	out := fmt.Sprintf("func Provide() (%s) {\n", returnedTypes)

	for _, expr := range anti.exprs {
		out += fmt.Sprintf("\t%s\n", expr)
	}
	out += "}\n"

	return out
}

func (anti *AntiDig) returnStmt(invokedType reflect.Type) string {
	returnStmt := "return "
	for i := 0; i < invokedType.NumIn(); i++ {
		typ := invokedType.In(i)
		returnStmt += fmt.Sprintf("%s, ", anti.Varname(typ))
	}
	return strings.TrimRight(returnStmt, ", ")
}

func (anti *AntiDig) returnedTypes(invokedType reflect.Type) string {
	returnStmt := ""
	for i := 0; i < invokedType.NumIn(); i++ {
		typ := invokedType.In(i)
		prefix := ""

		if typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
			prefix = "*"
		}
		returnStmt += fmt.Sprintf("%s%s.%s, ", prefix, anti.PkgAlias(typ.PkgPath()), typ.Name())
	}
	return strings.TrimRight(returnStmt, ", ")
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

func (anti *AntiDig) NamedVarname(name string, typ reflect.Type) string {
	key := typeAlias{Type: typ, Alias: name}

	varname, ok := anti.varnames[key]
	if ok {
		return varname
	}

	anti.varnameSeq++
	varname = fmt.Sprintf("var%d", anti.varnameSeq)
	if name != "" {
		varname += "_" + strcase.ToLowerCamel(name)
	}

	anti.varnames[key] = varname

	return varname
}

func (anti *AntiDig) Varname(typ reflect.Type) string {
	return anti.NamedVarname("", typ)
}

func (anti *AntiDig) GrouppedVarname(group string, typ reflect.Type) string {
	key := typeAlias{Type: typ, Alias: group}

	varname, ok := anti.groupVarnames[key]
	if ok {
		return varname
	}

	anti.varnameSeq++
	varname = fmt.Sprintf("var%d", anti.varnameSeq)
	if group != "" {
		varname += "_" + strcase.ToLowerCamel(group)
	}

	anti.groupVarnames[key] = varname

	return varname
}

func (anti *AntiDig) Seqname(typ reflect.Type) string {
	return anti.GrouppedSeqname("", typ)
}

func (anti *AntiDig) GrouppedSeqname(group string, typ reflect.Type) string {
	key := typeAlias{Type: typ, Alias: group}
	varname := anti.GrouppedVarname(group, typ)

	_, ok := anti.seqnames[key]
	if !ok {
		anti.seqnames[key] = 0
		return fmt.Sprintf("%s_0", varname)
	}

	anti.seqnames[key]++
	return fmt.Sprintf("%s_%d", varname, anti.seqnames[key])
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

	if anti.allPkgAliases[alias] {
		newalias := path[len(path)-2] + path[len(path)-1]

		for i := 0; anti.allPkgAliases[newalias]; i++ {
			newalias = fmt.Sprintf("%s%d", alias, i)
		}

		alias = newalias
	}

	anti.pkgAlias[pkgname] = alias
	anti.allPkgAliases[alias] = true

	return alias
}

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()
var errExpr = []string{"if err != nil {", "\tpanic(err)", "}"}

func (anti *AntiDig) SetErrorExpr(fntype reflect.Type) {
	errStmt := "\treturn "
	for i := 0; i < fntype.NumIn(); i++ {
		typ := fntype.In(i)

		switch typ.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Interface:
			errStmt += "nil, "
		default:
			alias := Anti.PkgAlias(typ.PkgPath())
			errStmt += fmt.Sprintf("%s.%s{}, ", alias, typ.Name())
		}
	}
	errStmt = strings.TrimRight(errStmt, ", ")
	errExpr = []string{"if err != nil {", errStmt, "}"}
}
