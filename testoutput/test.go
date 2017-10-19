package test

import (
	"github.com/kvigen/go-linq/template"
	"github.com/kvigen/go-linq/types"
)

func RunTheFilter(input []types.Input) []types.Output {
	output := template.Exec(template.SelectNode{})
	return output
}
