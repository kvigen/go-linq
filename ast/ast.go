package ast

import "reflect"

// TODO: support 'AS'
type SelectNode struct {
	// https://stackoverflow.com/questions/23030884/is-there-a-way-to-create-an-instance-of-a-struct-from-a-string
	InputType  string
	OutputType string
	// TODO: Support multiple fields
	Field string
	// TODO: In theory this should be any node. Nodes will have a simple interface
	// at some point...
	Src FromNode
}

// TODO: probably want to leave these in the generated code...
func (n *SelectNode) Next() interface{} {

	next := n.Src.Next()

	if next == nil {
		return nil
	}

	input := reflect.ValueOf(next)
	inputField := reflect.Indirect(input).FieldByName(n.Field)

	// TOOD: Fix the type if we decide to use this approach...
	output := reflect.New(reflect.TypeOf(SelectNode{}))
	outputField := reflect.Indirect(output).FieldByName(n.Field)

	outputField.Set(inputField.Elem())

	return output.Elem()
}

type FromNode struct {
	// TODO: Support reading from files and similar things...
	Data       []interface{}
	OutputType string
	current    int
}

func (n *FromNode) Next() interface{} {
	if n.current >= len(n.Data) {
		return nil
	}
	toReturn := n.Data[n.current]
	n.current++
	return toReturn
}
