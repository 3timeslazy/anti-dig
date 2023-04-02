// Copyright (c) 2021 Uber Technologies, Inc.
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
	"fmt"
	"reflect"
	"strings"

	"github.com/3timeslazy/anti-dig/internal/digerror"
	"github.com/3timeslazy/anti-dig/internal/digreflect"
	"github.com/3timeslazy/anti-dig/internal/dot"
)

// constructorNode is a node in the dependency graph that represents
// a constructor provided by the user.
//
// constructorNodes can produce zero or more values that they store into the container.
// For the Provide path, we verify that constructorNodes produce at least one value,
// otherwise the function will never be called.
type constructorNode struct {
	ctor  interface{}
	ctype reflect.Type

	// Location where this function was defined.
	location *digreflect.Func

	// id uniquely identifies the constructor that produces a node.
	id dot.CtorID

	// Whether the constructor owned by this node was already called.
	called bool

	// Type information about constructor parameters.
	paramList paramList

	// Type information about constructor results.
	resultList resultList

	// order of this node in each Scopes' graphHolders.
	orders map[*Scope]int

	// scope this node is part of
	s *Scope

	// scope this node was originally provided to.
	// This is different from s if and only if the constructor was Provided with ExportOption.
	origS *Scope
}

type constructorOptions struct {
	// If specified, all values produced by this constructor have the provided name
	// belong to the specified value group or implement any of the interfaces.
	ResultName  string
	ResultGroup string
	ResultAs    []interface{}
	Location    *digreflect.Func
}

func newConstructorNode(ctor interface{}, s *Scope, origS *Scope, opts constructorOptions) (*constructorNode, error) {
	cval := reflect.ValueOf(ctor)
	ctype := cval.Type()
	cptr := cval.Pointer()

	params, err := newParamList(ctype, s)
	if err != nil {
		return nil, err
	}

	results, err := newResultList(
		ctype,
		resultOptions{
			Name:  opts.ResultName,
			Group: opts.ResultGroup,
			As:    opts.ResultAs,
		},
	)
	if err != nil {
		return nil, err
	}

	location := opts.Location
	if location == nil {
		location = digreflect.InspectFunc(ctor)
	}

	n := &constructorNode{
		ctor:       ctor,
		ctype:      ctype,
		location:   location,
		id:         dot.CtorID(cptr),
		paramList:  params,
		resultList: results,
		orders:     make(map[*Scope]int),
		s:          s,
		origS:      origS,
	}
	s.newGraphNode(n, n.orders)
	return n, nil
}

func (n *constructorNode) Location() *digreflect.Func { return n.location }
func (n *constructorNode) ParamList() paramList       { return n.paramList }
func (n *constructorNode) ResultList() resultList     { return n.resultList }
func (n *constructorNode) ID() dot.CtorID             { return n.id }
func (n *constructorNode) CType() reflect.Type        { return n.ctype }
func (n *constructorNode) Order(s *Scope) int         { return n.orders[s] }
func (n *constructorNode) OrigScope() *Scope          { return n.origS }

func (n *constructorNode) String() string {
	return fmt.Sprintf("deps: %v, ctor: %v", n.paramList, n.ctype)
}

// Call calls this constructor if it hasn't already been called and
// injects any values produced by it into the provided container.
func (n *constructorNode) Call(c containerStore) (err error) {
	if n.called {
		return nil
	}

	if err := shallowCheckDependencies(c, n.paramList); err != nil {
		return errMissingDependencies{
			Func:   n.location,
			Reason: err,
		}
	}

	if n.s.recoverFromPanics {
		defer func() {
			if p := recover(); p != nil {
				err = PanicError{
					fn:    n.location,
					Panic: p,
				}
			}
		}()
	}

	args, err := n.paramList.BuildList(c)
	if err != nil {
		return errArgumentsFailed{
			Func:   n.location,
			Reason: err,
		}
	}

	loc := n.Location()
	resultsList := n.ResultList().Results

	argsVars := []string{}
	for _, param := range n.ParamList().Params {
		switch param := param.(type) {
		case paramSingle:
			paramVarname := Varname(param.Type)
			argsVars = append(argsVars, paramVarname)

		case paramObject:
			paramVarname := Varname(param.Type)
			argsVars = append(argsVars, paramVarname)

			if param.Type.Kind() != reflect.Struct {
				panic("only Struct params accepted")
			}

			beforeExpr := ""
			afterExpr := ""

			// Param's fields
			callexpr := fmt.Sprintf("%s := %s{\n", paramVarname, param.Type)
			for _, field := range param.Fields {
				switch param := field.Param.(type) {
				case paramSingle:
					fieldVarname := Varname(param.Type)
					callexpr += fmt.Sprintf("\t\t%s: %s,\n", field.FieldName, fieldVarname)

				case paramGroupedSlice:
					fieldVarname := Varname(param.Type)
					callexpr += fmt.Sprintf("\t\t%s: %s,\n", field.FieldName, fieldVarname)

					// Create the slice before param
					beforeExpr += fmt.Sprintf("%s := %s{\n", fieldVarname, param.Type)

					elem := param.Type.Elem()
					PkgAlias(elem.PkgPath())

					suffix := ""
					if path, ok := TypeWrapper[elem]; ok {
						suffix = path
					}

					cnt := TypeSlicename[elem]
					elemVarname := Varname(elem)
					for i := 0; i <= cnt; i++ {
						varname := fmt.Sprintf("%s_%d", elemVarname, i)
						if FlattenVars[varname] {
							suffix := FlattenVarsSuffixies[varname]
							afterExpr += fmt.Sprintf("%s = append(%s, %s.%s...)\n", fieldVarname, fieldVarname, varname, suffix)
							continue
						}
						beforeExpr += fmt.Sprintf("\t\t%s.%s,\n", varname, suffix)
					}

					beforeExpr += "\t}"

				default:
					panic("Unsupported field type")
				}
			}
			callexpr += "\t}"
			CallExprs = append(CallExprs, beforeExpr)
			CallExprs = append(CallExprs, afterExpr)
			CallExprs = append(CallExprs, callexpr)

		default:
			panic("govno")
		}
	}

	returnsErr := ErrorFunc(n.ResultList())

	switch returned := resultsList[0].(type) {
	case resultSingle:
		varname := Varname(returned.Type)

		funcexpr := fmt.Sprintf("%s.%s", PkgAlias(loc.Package), loc.Name)
		expr := fmt.Sprintf("%s := %s(%s)", varname, funcexpr, strings.Join(argsVars, ", "))
		if returnsErr {
			expr = fmt.Sprintf("%s, err := %s(%s)", varname, funcexpr, strings.Join(argsVars, ", "))
			expr += "\n" + PanicOnErr
		}

		CallExprs = append(CallExprs, expr)

	case resultObject: // dig.Out
		// TODO: many
		for _, field := range returned.Fields {
			switch result := field.Result.(type) {
			case resultGrouped:
				TypeWrapper[result.Type] = returned.Type.Field(1).Name

				varname := Slicename(result.Type)
				if result.Flatten {
					FlattenVars[varname] = result.Flatten
					FlattenVarsSuffixies[varname] = returned.Type.Field(1).Name
				}

				funcexpr := fmt.Sprintf("%s.%s", PkgAlias(loc.Package), loc.Name)
				expr := fmt.Sprintf("%s := %s(%s)", varname, funcexpr, strings.Join(argsVars, ", "))
				if returnsErr {
					expr = fmt.Sprintf("%s, err := %s(%s)", varname, funcexpr, strings.Join(argsVars, ", "))
					expr += "\n" + PanicOnErr
				}

				CallExprs = append(CallExprs, expr)

			default:
				panic("field.Result.(type)")
			}
		}

	default:
		panic("govno")
	}

	receiver := newStagingContainerWriter()
	results := c.invoker()(reflect.ValueOf(n.ctor), args)
	if err := n.resultList.ExtractList(receiver, false /* decorating */, results); err != nil {
		return errConstructorFailed{Func: n.location, Reason: err}
	}

	// Commit the result to the original container that this constructor
	// was supplied to. The provided constructor is only used for a view of
	// the rest of the graph to instantiate the dependencies of this
	// container.
	receiver.Commit(n.s)
	n.called = true

	return nil
}

// stagingContainerWriter is a containerWriter that records the changes that
// would be made to a containerWriter and defers them until Commit is called.
type stagingContainerWriter struct {
	values map[key]reflect.Value
	groups map[key][]reflect.Value
}

var _ containerWriter = (*stagingContainerWriter)(nil)

func newStagingContainerWriter() *stagingContainerWriter {
	return &stagingContainerWriter{
		values: make(map[key]reflect.Value),
		groups: make(map[key][]reflect.Value),
	}
}

func (sr *stagingContainerWriter) setValue(name string, t reflect.Type, v reflect.Value) {
	sr.values[key{t: t, name: name}] = v
}

func (sr *stagingContainerWriter) setDecoratedValue(_ string, _ reflect.Type, _ reflect.Value) {
	digerror.BugPanicf("stagingContainerWriter.setDecoratedValue must never be called")
}

func (sr *stagingContainerWriter) submitGroupedValue(group string, t reflect.Type, v reflect.Value) {
	k := key{t: t, group: group}
	sr.groups[k] = append(sr.groups[k], v)
}

func (sr *stagingContainerWriter) submitDecoratedGroupedValue(_ string, _ reflect.Type, _ reflect.Value) {
	digerror.BugPanicf("stagingContainerWriter.submitDecoratedGroupedValue must never be called")
}

// Commit commits the received results to the provided containerWriter.
func (sr *stagingContainerWriter) Commit(cw containerWriter) {
	for k, v := range sr.values {
		cw.setValue(k.name, k.t, v)
	}

	for k, vs := range sr.groups {
		for _, v := range vs {
			cw.submitGroupedValue(k.group, k.t, v)
		}
	}
}
