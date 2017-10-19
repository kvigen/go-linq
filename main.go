package main

import (
	"fmt"

	// TODO: rename this???
	testoutput "github.com/kvigen/go-linq/testoutput"
	"github.com/kvigen/go-linq/types"
)

func main() {
	input2 := []types.Input{types.Input{"a", "b"}, types.Input{"c", "d"}}
	output := testoutput.RunTheFilter(input2)
	fmt.Printf("Output: %+v\n", output)
}

func Exec(sql string) interface{} {
	return nil
	// TODO: Convert it to an AST and execute that
}

// Future fields
// SELECT input1 AS output1
// JOINs
// GROUP BY
// LIMIT
