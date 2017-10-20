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

	inputs2 := []interface{}{input{"xa"}}

	otherTable := From(inputs2)
	// results := From(inputs).Map(func(str string) string { return str + "b" }).Array()
	results := From(inputs).
		Map(func(i input) intermediate { return intermediate{Str: i.Str + "a"} }).
		Join(otherTable, func(i intermediate) string { return i.Str }, func(i input) string { return i.Str }).
		//Map(func(i intermediate) output { return output{Str: i.Str + "b"} }).
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

func (t Table) Join(t2 Table, f1 interface{}, f2 interface{}) Table {
	// TODO: Add some form of caching here...

	return Table{iTable: &joinTable{itLeft: t, itRight: t2, idLeft: reflect.ValueOf(f1), idRight: reflect.ValueOf(f2)}}
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

// TODO: support different kinds of joins...
type joinTable struct {
	itRight Table
	itLeft  Table
	idRight reflect.Value // function
	idLeft  reflect.Value // function
	right   map[interface{}]interface{}
}

func (t *joinTable) Next() (interface{}, bool) {
	if t.right == nil {
		// Setup the hash table
		t.right = make(map[interface{}]interface{})
		for {
			res, ok := t.itRight.Next()
			if !ok {
				break
			}
			input := []reflect.Value{reflect.ValueOf(res)}
			output := t.idRight.Call(input)
			o := output[0].Interface()
			// TODO: this assumes that only one things has a given key
			t.right[o] = res
		}
	}

	res, ok := t.itLeft.Next()
	if !ok {
		return nil, false
	}
	input := []reflect.Value{reflect.ValueOf(res)}
	output := t.idLeft.Call(input)
	o := output[0].Interface()
	// TODO: Right now we're assuming that there's only one match
	val := t.right[o]
	// TODO: name these variables better...
	return map[string]interface{}{
		"lhs": res,
		"rhs": val,
	}, true
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
