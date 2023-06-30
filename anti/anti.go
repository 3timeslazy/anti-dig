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

package anti

import (
	"container/list"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"reflect"
	"strings"

	"github.com/3timeslazy/anti-dig/anti/optimise"
	"github.com/3timeslazy/anti-dig/anti/overwrite"

	"github.com/iancoleman/strcase"
	"golang.org/x/exp/slices"
)

type AntiDig struct {
	output   io.Writer
	exprs    []string
	optimise bool
	rename   bool

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

	varTypes map[string]string
}

type typeAlias struct {
	Alias string
	Type  reflect.Type
}

func New(output io.Writer) *AntiDig {
	return &AntiDig{
		output:       output,
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

		varTypes: map[string]string{},
	}
}

func (anti *AntiDig) Optimise(enable bool) *AntiDig {
	anti.optimise = enable
	return anti
}

func (anti *AntiDig) Rename(enable bool) *AntiDig {
	anti.rename = enable
	return anti
}

func (anti *AntiDig) Generate(invokedType reflect.Type) error {
	generated := strings.Join([]string{
		"package main\n",
		anti.generateImports(),
		anti.generateFunc(invokedType),
	}, "\n")

	fset := token.NewFileSet()
	generatedFile, err := parser.ParseFile(fset, "provide.gen.go", generated, parser.ParseComments)
	if err != nil {
		return err
	}

	if anti.rename {
		err = overwrite.Rename(fset, generatedFile)
		if err != nil {
			return err
		}
	}
	if anti.optimise {
		fset, generatedFile, err = optimise.Optimise(fset, generatedFile, anti.varTypes)
		if err != nil {
			return err
		}
	}

	return format.Node(anti.output, fset, generatedFile)
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
		alias := typeAlias{Type: typ}
		varname := ""

		if _, ok := anti.varnames[alias]; ok {
			varname = anti.Varname(typ)
		} else {
			varname = anti.groupVarnames[alias]
			varname += fmt.Sprintf("_%d", anti.seqnames[alias])
		}

		returnStmt += fmt.Sprintf("%s, ", varname)
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
	pkgs := []string{}

	for pkg, alias := range anti.pkgAlias {
		pkgs = append(pkgs, fmt.Sprintf("\t%s \"%s\"", alias, pkg))
	}

	// Sort the imports so that their order is always
	// the same and there are no fluctuations in the tests associated with this
	slices.SortFunc(pkgs, func(a, b string) bool {
		return a < b
	})

	out := "import (\n" + strings.Join(pkgs, "\n")
	out = out + "\n)\n"

	return out
}

func (anti *AntiDig) PushFnCall(fnName string) {
	anti.callstack.PushBack(fnName)
}

func (anti *AntiDig) PopFnCall() {
	anti.cleanupFn(anti.currFn())
	elem := anti.callstack.Back()
	anti.callstack.Remove(elem)
}

func (anti *AntiDig) cleanupFn(fnName string) {
	delete(anti.fnsArgs, fnName)
	delete(anti.fnsVars, fnName)
	delete(anti.fnsSuffixies, fnName)
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

var errExpr = []string{"if err != nil {", "\tpanic(err)", "}"}

func (anti *AntiDig) SetErrorExpr(fntype reflect.Type) {
	errStmt := "\treturn "
	for i := 0; i < fntype.NumIn(); i++ {
		typ := fntype.In(i)

		switch typ.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice, reflect.Interface:
			errStmt += "nil, "
		default:
			alias := anti.PkgAlias(typ.PkgPath())
			errStmt += fmt.Sprintf("%s.%s{}, ", alias, typ.Name())
		}
	}
	errStmt = strings.TrimRight(errStmt, ", ")
	errExpr = []string{"if err != nil {", errStmt, "}"}
}

func (anti *AntiDig) PrintErrorExpr() {
	anti.Print(errExpr...)
}

func (anti *AntiDig) AddVarType(varname string, typ reflect.Type) {
	typestr := typ.Name()
	if typestr == "" {
		typestr = typ.Elem().Name()
	}

	anti.varTypes[varname] = typestr
}
