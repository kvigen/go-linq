package ast

import "strings"

func Parse(sql, inputType, outputType string) SelectNode {
	// Handles two query forms:
	// SELECT x FROM a
	// SELECT x FROM a JOIN b ON a = b
	// TODO: Make this much smarter...
	selectPos := strings.Index(sql, "SELECT")
	fromPos := strings.Index(sql, "FROM")
	joinPos := strings.Index(sql, "JOIN")
	onPos := strings.Index(sql, "ON")
	if selectPos == -1 || fromPos == -1 {
		panic("invalid SQL -- this should really be an error not a panic")
	}
	if fromPos < selectPos {
		panic("invalid SQL -- this should really be an error not a panic")
	}
	if joinPos != -1 {
		if onPos == -1 {
			panic("invalid SQL -- this should really be an error not a panic")
		}
		if joinPos < fromPos || onPos < joinPos {
			panic("invalid SQL -- this should really be an error not a panic")
		}
	}

	// TODO: Bounds checking... or really just a better lexer / parser

	if joinPos == -1 {
		selectStatement := sql[selectPos+7 : fromPos]
		fromStatement := sql[fromPos+5 : len(sql)-1]
		fromNode := FromNode{
			OutputType: inputType,
			Variable:   fromStatement,
		}
		selectNode := SelectNode{
			InputType:  inputType,
			OutputType: outputType,
			Field:      selectStatement,
			Src:        fromNode,
		}
		return selectNode
	} else {
		selectStatement := sql[selectPos+7 : fromPos]
		from1Statement := sql[fromPos+5 : joinPos]
		from2Statement := sql[joinPos+5 : onPos+2]
		onStatement := sql[onPos+2:]

		from1Node := FromNode{
			OutputType: inputType,
			Variable:   from1Statement,
		}
		from2Node := FromNode{
			OutputType: inputType,
			Variable:   from2Statement,
		}
		joinNode := JoinNode{
			LeftSide:  from1Node,
			RightSide: from2Node,
			On:        onStatement,
		}
		selectNode := SelectNode{
			InputType:  inputType,
			OutputType: outputType,
			Field:      selectStatement,
			Src:        joinNode,
		}
	}

}

// TODO: support 'AS'
type SelectNode struct {
	// https://stackoverflow.com/questions/23030884/is-there-a-way-to-create-an-instance-of-a-struct-from-a-string
	InputType  string
	OutputType string
	// TODO: Support multiple fields
	Field string
	// TODO: In theory this should be any node. Nodes will have a simple interface
	// at some point...
	Src interface{}
}

type FromNode struct {
	// TODO: Support reading from files and similar things...
	OutputType string
	Variable   string
}

type JoinNode struct {
	LeftSide  FromNode
	RightSide FromNode
	On        string
}

/*
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


func (n *FromNode) Next() interface{} {
	if n.current >= len(n.Data) {
		return nil
	}
	toReturn := n.Data[n.current]
	n.current++
	return toReturn
}
*/
