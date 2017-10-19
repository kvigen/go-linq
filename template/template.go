
package template

// TODO: This is a major hack... though we can require people to pass in the types
// use or we can figure it out from the types passed in... that sounds pretty good...
import (
	"github.com/kvigen/go-linq/types"
)



type SelectNode struct {
	It *FromNode
}

func (s *SelectNode) Next() *types.Output {

	for {
		next := s.It.Next()
		if next == nil {
			return nil
		}
		o := types.Output{next.Field1}
		return &o
	}
}


type FromNode struct {
	Data    []types.Input
	current int
}	

func (f *FromNode) Next() *types.Input {
	if f.current >= len(f.Data) {
		return nil
	}
	toReturn := f.Data[f.current]
	f.current++
	return &toReturn
}


func Exec(node SelectNode) []types.Output{

	results := make([]types.Output, 0)

	for {
		result := node.Next()
		if result == nil {
			break
		}
		results = append(results, *result)
	}
	return results
}


