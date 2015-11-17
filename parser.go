package main

import (
	"fmt"
	"go/parser"
	"go/token"
)

func mainOld() {
	fset := token.NewFileSet() // positions are relative to fset

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, "input.go", nil, parser.Trace)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the imports from the file's AST.
	for _, s := range f.Decls {
		fmt.Println(s) //Some func
		// fmt.Println(s.Pos()) //Some numbers
		// fmt.Println(fset.Position(s.Pos()) //Some numbers
		// fmt.Println(fset.Position(s.End())) //Some numbers

	}

}
