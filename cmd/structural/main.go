package main

import (
	"fmt"
	"reflect"
)

type input struct {
	Str string
}

type intermediate struct {
	Str string
}

type output struct {
	Str string
}

func main() {

	inputs := []interface{}{input{"x"}, input{"y"}}

	//results := From(inputs).Map(func(str string) string { return str + "a" }).
	//	Map(func(str string) string { return str + "b" }).Array()
	// TODO: Verify that the below is wrong...
	// results := From(inputs).Map(func(str string) string { return str + "b" }).Array()
	results := From(inputs).
		Map(func(i input) intermediate { return intermediate{Str: i.Str + "a"} }).
		Map(func(i intermediate) output { return output{Str: i.Str + "b"} }).
		Array()
	fmt.Printf("Result %+v\n", results)

}

type tableInterface interface {
	// TODO: Do we actually need the bool val? Nil check can be tricky...
	Next() (interface{}, bool)
}

type Table struct {
	// composition
	iTable tableInterface
}

func (t Table) Next() (interface{}, bool) {
	return t.iTable.Next()
}

// TODO: Figure out how we might make this generic...
func From(input []interface{}) Table {
	return Table{iTable: &fromTable{data: input}}
}

// TODO: Can you have headless map functions so they can be chained independently...
func (t Table) Map(fn interface{}) Table {

	// TODO: special case for func(map[string]interface{}) map[string]interface{} since very common
	// This is making lots of assumptions that people aren't passing in crazy things...
	typ := reflect.TypeOf(fn)
	inType := typ.In(0)
	outType := typ.Out(0)

	if inType == reflect.TypeOf(intermediate{}) && outType == reflect.TypeOf(output{}) {
		cacheMap := cacheMapFunc{it: t, fn: fn.(func(intermediate) output)}
		return Table{iTable: &cacheMap}
	}
	// just a different return object...
	return Table{iTable: &mapTable{it: t, fnValue: reflect.ValueOf(fn)}}
}

func (t Table) Array() []interface{} {
	results := make([]interface{}, 0)
	for {
		n, ok := t.Next()
		if !ok {
			break
		}
		results = append(results, n)
	}
	return results
}

type fromTable struct {
	data         []interface{}
	currentIndex int
}

func (t *fromTable) Next() (interface{}, bool) {
	if t.currentIndex >= len(t.data) {
		return nil, false
	}
	ret := t.data[t.currentIndex]
	t.currentIndex++
	return ret, true
}

type mapTable struct {
	fnValue reflect.Value
	it      Table
}

func (t *mapTable) Next() (interface{}, bool) {
	val, ok := t.it.Next()
	if !ok {
		return nil, false
	}

	input := []reflect.Value{reflect.ValueOf(val)}
	output := t.fnValue.Call(input)
	return output[0].Interface(), true
}

// in general this would be auto-generated. A generics wrapper...
type cacheMapFunc struct {
	it Table
	fn func(intermediate) output
}

func (t *cacheMapFunc) Next() (interface{}, bool) {
	fmt.Printf("in cached map\n")
	res, ok := t.it.Next()
	if !ok {
		return nil, ok
	}

	i := res.(intermediate)
	return t.fn(i), true
}
