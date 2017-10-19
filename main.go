package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kvigen/go-linq/ast"
	"github.com/kvigen/go-linq/generator"
	"github.com/kvigen/go-linq/template"
	// TODO: rename this???
	"github.com/kvigen/go-linq/transform"
	"github.com/kvigen/go-linq/types"
)

func main() {

	action := flag.String("action", "", "generate or run or transform")
	flag.Parse()

	// TODO: Support something other than []interface{}. Not sure how
	// to pass this information into FromNode that way? Maybe that could
	// just take in an interface?
	input := []interface{}{types.Input{"a", "b"}, types.Input{"c", "d"}}

	fromNode := ast.FromNode{
		OutputType: "types.Input",
		Data:       input,
	}

	selectNode := ast.SelectNode{
		InputType:  "types.Input",
		OutputType: "types.Output",
		Field:      "Field1",
		Src:        fromNode,
	}

	if *action == "generate" {

		//output := Exec("SELECT Output{field1} FROM input").(Output)
		//output := executor.Exec(selectNode)
		outputCode, err := generator.Generate(selectNode)
		if err != nil {
			log.Fatalf("Failed to generate code: %s\n", err)
		}
		if err := ioutil.WriteFile("template/template.go", []byte(outputCode), 0644); err != nil {
			log.Fatalf("Error writing template: %s\n", err)
		}
	} else if *action == "transform" {
		srcBytes, err := ioutil.ReadFile("testinput/test.go")
		if err != nil {
			log.Fatalf("Error reading transform file: %s\n", err)
		}
		outputCode, _, err := transform.TransformFile(string(srcBytes))
		if err != nil {
			log.Fatalf("Error transforming file: %s", err)
		}
		if err := ioutil.WriteFile("testoutput/test.go", []byte(outputCode), 0644); err != nil {
			log.Fatalf("Error writing template: %s\n", err)
		}

	} else {
		fromNode := template.FromNode{Data: input}
		selectNode := template.SelectNode{It: &fromNode}
		output := template.Exec(selectNode)
		fmt.Printf("Output: %+v\n", output)
	}

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
