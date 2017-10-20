package linq

// Run is a placeholder for the actual generated code
func Run(interface{}) interface{} {
	return nil
}

// TODO: is this the use case for the concurrent cache??? this should be initialized by
// the pre-compiler
// TODO: Should this map hold whole tables... not sure if that has anything we need...
var readCache map[string]tableIt = make(map[string]tableIt)

type Table struct {
	// specific implementation??? this part should be an interface... I think
	// TODO: is this optional???
	it tableIt
}

// TODO: Think about whether this is the best interface
func (t Table) HasNext() bool {
	return t.it.HasNext()
}

func (t Table) Next() interface{} {
	return t.it.Next()
}

func (t Table) Select(field string) Table {
	key := t.it.OutputType() + ":" + field
	tableIt, ok := readCache[key]
	if !ok {
		panic("should be able to find this...")
	}
	return Table{it: tableIt}
}

// implementation of table -- using composition to simluate inheritance
type tableIt interface {
	HasNext() bool
	Next() interface{}
	OutputType() string
}

func From(input interface{}) Table {
	// If type of input is something we've pre-computed
	return precomputedTable
}
