package transform

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"strings"
)

// TODO: These should obviously not be global variables...
var start int
var end int

// TransformFile returns the new file as a string and a bool to indicate
// if it actually did any transformations
func TransformFile(src string) (string, bool, error) {

	fset := token.NewFileSet()

	parsed, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return "", false, err
	}

	// Need to fix this to actually find and replace the linq import
	parsed.Imports[0].Path.Value = "\"github.com/kvigen/go-linq/template\""
	/*

	 */

	ast.Walk(visitor{}, parsed)

	/*
		for _, decl := range parsed.Decls {
			fnDecl, ok := decl.(*ast.FuncDecl)
			if ok {

				// Go through the body and process...
				// TODO: Support going deeper into all the types with bodies...
				// Shouldn't be too hard to build a recursive function to do this...
				for _, stmt := range fnDecl.Body.List {
					if assign, ok := stmt.(*ast.AssignStmt); ok {
						for _, expr := range assign.Rhs {
							fmt.Printf("RHS: %+v\n", expr)
							if callExpr, ok := expr.(*ast.CallExpr); ok {
								if fmt.Sprintf("%s", callExpr.Fun.(*ast.SelectorExpr).X) == "linq" {
									fmt.Println("Found it!!!")
								}
							}
						}
					}
				}
			}
		}
	*/

	var buf bytes.Buffer
	printer.Fprint(&buf, fset, parsed)
	str := buf.String()

	log.Printf("Start: %d, End: %d, Total: %d\n", start, end, len(str))

	// TODO: hard-coding these is obviously wrong... In particular it doesn't
	// handle comments correctly... Somehow doesn't mix nicely with the imports...
	result := str[:start+3]
	result += "template.Exec(template.SelectNode{})"
	result += str[end+24:]
	fmt.Println(result)

	return result, false, nil
}

type visitor struct {
}

type parsedSQL struct {
	Select string
	From   string
}

func parse(sql string) parsedSQL {
	// TOOD: Make this much more sophisticated

	// Find "SELECT" AND "FROM"
	// TODO: Make this much smarter...
	selectPos := strings.Index(sql, "SELECT")
	fromPos := strings.Index(sql, "FROM")
	if selectPos == -1 || fromPos == -1 || fromPos < selectPos {
		panic("invalid SQL -- this should really be an error not a panic")
	}

	// TODO: Bounds checking... or really just a better lexer / parser
	selectStatement := sql[selectPos+7 : fromPos]
	fromStatement := sql[fromPos+5:]

	return parsedSQL{Select: selectStatement, From: fromStatement}
}

func (v visitor) Visit(node ast.Node) ast.Visitor {

	// TODO: This isn't working...

	/*
		if file, ok := node.(*ast.File); ok {
			file.Imports = append(file.Imports, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "github.com/kvigen/go-linq/template",
				},
			})
		}
	*/

	// TODO: Also look for type assertions so we can set the output type
	// in it to use later...

	if callExpr, ok := node.(*ast.CallExpr); ok {

		start = int(callExpr.Pos())
		end = int(callExpr.End())

		expr := callExpr.Fun.(*ast.SelectorExpr).X
		if fmt.Sprintf("%s", expr) == "linq" {
			log.Printf("Found it %T!!!\n", expr)
		}
		if len(callExpr.Args) == 0 {
			panic("something has gone horribly wrong")
		}
		fmt.Printf("Type: %T\n", callExpr.Args[0])
		arg, ok := callExpr.Args[0].(*ast.BasicLit)
		if !ok {
			panic("something has gone horribly wrong")
		}
		sql := parse(arg.Value)
		fmt.Printf("SQL: %+v\n", sql)

		iden := expr.(*ast.Ident)
		iden.Name = "template"
	}

	return v
}
