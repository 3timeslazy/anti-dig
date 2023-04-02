package antidig

import (
	"fmt"
	"reflect"
	"strings"
)

var FlattenVars = map[string]bool{}
var FlattenVarsSuffixies = map[string]string{}

var TypeWrapper = map[reflect.Type]string{}

var TypeSlicename = map[reflect.Type]int{}

func Slicename(typ reflect.Type) string {
	varname := Varname(typ)

	_, ok := TypeSlicename[typ]
	if !ok {
		TypeSlicename[typ] = 0
		return fmt.Sprintf("%s_0", varname)
	}

	TypeSlicename[typ] += 1
	return fmt.Sprintf("%s_%d", varname, TypeSlicename[typ])
}

var TypeVarname = map[reflect.Type]string{}
var Seq int

func Varname(typ reflect.Type) string {
	varname, ok := TypeVarname[typ]
	if ok {
		return varname
	}

	Seq++
	TypeVarname[typ] = fmt.Sprintf("var%d", Seq)
	return TypeVarname[typ]
}

var PackageAlias = map[string]string{}
var PkgSeq int

func PkgAlias(pkgname string) string {
	alias, ok := PackageAlias[pkgname]
	if ok {
		return alias
	}

	path := strings.Split(pkgname, "/")
	alias = path[len(path)-1]
	PackageAlias[pkgname] = alias
	return alias

	// PkgSeq++
	// path := strings.Split(pkgname, "/")

	// alias = fmt.Sprintf("%s%d", path[len(path)-1], PkgSeq)
	// PackageAlias[pkgname] = alias
	// return alias
}

var CallExprs []string

func ErrorFunc(results resultList) bool {
	for _, idx := range results.resultIndexes {
		if idx == -1 {
			return true
		}
	}

	return false
}

var PanicOnErr = "\tif err != nil {\n\t\tpanic(err)\n\t}"

func PrintGenerated() {
	fmt.Println("package main")
	fmt.Println()

	fmt.Println("import (")
	for pkg, alias := range PackageAlias {
		fmt.Printf("\t%s \"%s\"\n", alias, pkg)
	}
	fmt.Println(")")
	fmt.Println()

	fmt.Println("func main() {")
	for _, callexpr := range CallExprs {
		fmt.Printf("\t%s\n", callexpr)
	}
	fmt.Println("}")
}
