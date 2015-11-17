package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func oldMain() {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "", `
package p
import "fmt"
//
var someField int
//
func someFunc() {
	methodVar := "here"
	fmt.Println(methodVar)
}
`, 0)
	//if err != nil {
	//	t.Fatal(err)
	//}

	// Parse the file containing this very example
	// but stop after processing the imports.
	//f, err := parser.ParseFile(fset, "example_test.go", nil, parser.ImportsOnly)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the imports from the file's AST.
	for _, s := range f.Imports {
		fmt.Println(s.Path.Value)
	}

	fmt.Printf("f.Decls: %q\n", f.Decls)

	for _, d := range f.Decls {
		//Prints import and variable Declarations
		if d, ok := d.(*ast.GenDecl); ok {

			fmt.Printf("d: %q ok: %q\n", d, ok)
			for _, s := range d.Specs {
				if value, ok := s.(*ast.ValueSpec); ok {
					fmt.Printf("\tvariable: %q\n", value.Names[0])
				}
			}
		}
		//Function declarations
		if fu, ok := d.(*ast.FuncDecl); ok {
			fmt.Printf("fu: %q\n", fu)
			fmt.Printf("fu.Name: %q\n", fu.Name)
			fmt.Printf("fu.Body: %q\n", fu.Body)

			fmt.Printf("fu.Body.List: %q\n", fu.Body.List)
			for _, stmt := range fu.Body.List {
				fmt.Printf("stmt: %q\n", stmt)
				if assignSt, ok := stmt.(*ast.AssignStmt); ok {
					fmt.Printf("\tassignSt.Lhs[0]: %q\n", assignSt.Lhs[0])
				}
			}

		}
	}
}

//func getField(file *ast.File, fieldname string) *ast.Field {
// ...
//}
