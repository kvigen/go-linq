package executor

import (
	"github.com/kvigen/go-linq/ast"
)

// TODO: Does this have to take in a SelectNode?
func Exec(ast ast.SelectNode) []interface{} {
	// TODO: Would be great if we could return an interface of the right type...
	// this might be somewhat difficult???
	results := make([]interface{}, 0)

	for {
		result := ast.Next()
		if result == nil {
			break
		}
		results = append(results, result)
	}
	return results
}
