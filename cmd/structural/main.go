package main

import "github.com/kvigen/go-linq/types"

func main() {

	input := []types.Input{types.Input{"a", "c"}, types.Input{"b", "d"}}

	linq.From(input).Select(input.Field1)

	From(input).SelectFunc(func() {})

}

// TODO: Figure out how we might make this generic...
func From(input []interface{}) FromTable {

}

type FromTable struct {
}

func Next() interface{} {

}

type Table interface {
	Next() interface{}
}

func (t Table) Select(field) Table {
	return SelectTable{t}
}

type SelectTable struct {
	T Table
}

func (t SelectTable) Next() interface{} {
	val := t.Next()
	// Use reflection... is there any way to remember the code here???
}
