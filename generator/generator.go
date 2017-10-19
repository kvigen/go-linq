package generator

import (
	"github.com/Clever/wag/templates"
	"github.com/kvigen/go-linq/ast"
)

// TODO: This should go in the generated code...
type iterator interface {
	Next() interface{}
}

type selectStruct struct {
	InputType  string
	OutputType string
	// TODO: support 'AS'
	Field string
}

// TODO: The select node is generic, so we shouldn't have to generate it...
var selectTemplate = `

type SelectNode struct {
	It *FromNode
}

func (s *SelectNode) Next() *{{.OutputType}} {

	for {
		next := s.It.Next()
		if next == nil {
			return nil
		}
		o := {{.OutputType}}{next.{{.Field}}}
		return &o
	}
}
`

type fromStruct struct {
	OutputType string
}

// Note that this is the from template for an array from clause... would want to support
// others...
var fromTemplate = `
type FromNode struct {
	Data    []interface{}
	current int
}	

func (f *FromNode) Next() *{{.OutputType}} {
	if f.current >= len(f.Data) {
		return nil
	}
	// TODO: Fix this...
	toReturn := f.Data[f.current].(types.Input)
	f.current++
	return &toReturn
}
`

type baseStruct struct {
	OutputType string
}

var baseTemplate = `
func Exec(node SelectNode) []{{.OutputType}}{

	results := make([]{{.OutputType}}, 0)

	for {
		result := node.Next()
		if result == nil {
			break
		}
		results = append(results, *result)
	}
	return results
}
`

type fullCodeStruct struct {
	Select string
	From   string
	Base   string
}

// TODO: I have
var fullCodeTemplate = `
package template

// TODO: This is a major hack... though we can require people to pass in the types
// use or we can figure it out from the types passed in... that sounds pretty good...
import (
	"github.com/kvigen/go-linq/types"
)

{{.Select}}
{{.From}}
{{.Base}}

`

func Generate(node ast.SelectNode) (string, error) {
	// Convert the ast nodes into execution nodes??? and return that too?

	sel := selectStruct{
		InputType:  node.InputType,
		OutputType: node.OutputType,
		Field:      node.Field,
	}
	selectCode, err := templates.WriteTemplate(selectTemplate, sel)
	if err != nil {
		return "", err
	}

	from := fromStruct{
		OutputType: node.Src.OutputType,
	}
	fromCode, err := templates.WriteTemplate(fromTemplate, from)
	if err != nil {
		return "", err
	}

	base := baseStruct{
		OutputType: node.OutputType,
	}
	baseCode, err := templates.WriteTemplate(baseTemplate, base)

	full := fullCodeStruct{
		Select: selectCode,
		From:   fromCode,
		Base:   baseCode,
	}

	return templates.WriteTemplate(fullCodeTemplate, full)
}