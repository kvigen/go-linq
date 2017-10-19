package test

import (
	"github.com/kvigen/go-linq/linq"
	"github.com/kvigen/go-linq/types"
)

func RunTheFilter(input []types.Input) []types.Output {
	output := linq.Run("SELECT field1 FROM input").([]types.Output)
	return output
}
