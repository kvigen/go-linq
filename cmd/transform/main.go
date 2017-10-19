package main

import (
	"io/ioutil"
	"log"

	"github.com/kvigen/go-linq/ast"
	"github.com/kvigen/go-linq/generator"
	"github.com/kvigen/go-linq/transform"
	"github.com/kvigen/go-linq/types"
)

func main() {

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
	// Really this should come in the parsing...not here

	//output := Exec("SELECT Output{field1} FROM input").(Output)
	//output := executor.Exec(selectNode)
	outputCode, err := generator.Generate(selectNode)
	if err != nil {
		log.Fatalf("Failed to generate code: %s\n", err)
	}
	if err := ioutil.WriteFile("../../template/template.go", []byte(outputCode), 0644); err != nil {
		log.Fatalf("Error writing template: %s\n", err)
	}

	srcBytes, err := ioutil.ReadFile("../../testinput/test.go")
	if err != nil {
		log.Fatalf("Error reading transform file: %s\n", err)
	}
	outputCode, _, err = transform.TransformFile(string(srcBytes))
	if err != nil {
		log.Fatalf("Error transforming file: %s", err)
	}
	if err := ioutil.WriteFile("../../testoutput/test.go", []byte(outputCode), 0644); err != nil {
		log.Fatalf("Error writing template: %s\n", err)
	}

}
