package antidig

import (
	"bytes"
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
	decls := [][]byte{
		[]byte("package main\n"),
		anti.generateImports(),
		anti.generateFunc(),
	}

	for _, decl := range decls {
		fmt.Fprintln(anti.output, string(decl))
	}
}

func (anti *AntiDig) generateFunc() []byte {
	buf := &bytes.Buffer{}

	fmt.Fprintln(buf, "func main() {")
	for _, expr := range anti.exprs {
		fmt.Fprintf(buf, "\t%s\n", expr)
	}
	fmt.Fprintln(buf, "}")

	return buf.Bytes()
}

func (anti *AntiDig) generateImports() []byte {
	buf := &bytes.Buffer{}

	fmt.Fprintln(buf, "import (")

	for pkg, alias := range anti.pkgAlias {
		fmt.Fprintf(buf, "\t%s \"%s\"\n", alias, pkg)
	}

	fmt.Fprintln(buf, ")")

	return buf.Bytes()
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

var PanicExpr = []string{"if err != nil {", "\tpanic(err)", "}"}
var ErrorInterface = reflect.TypeOf((*error)(nil)).Elem()
