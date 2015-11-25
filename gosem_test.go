package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"leonmoll.de/testutils"
)

// INFO: Print tree with:  ast.Print(fset, ...)

var originalArgs []string

func setup() {
	originalArgs = os.Args
	fmtPrintf = testutils.MockPrintf
	functionVarsFound = nil
	fieldsFound = nil
	functionsFound = nil
}

func teardown() {
	os.Args = originalArgs
	testutils.ResetMockPrinter()
}

func Test_main(t *testing.T) {
	defer teardown()
	setup()
	addToArgs("-f=input.go")
	expected := "someFieldToPrint someOtherField|8,16,stringToPrint a b|18,21,a"

	main()

	assertStringEqual(t, "Should have parsed input.go correctly", expected, testutils.GetLastPrinted())
}

func Test_main_shouldPrintErrorOnPanic(t *testing.T) {
	defer teardown()
	setup()
	addToArgs("-f=error.go")
	fmt.Println(os.Args)
	expected := ""

	main()

	assertStringEqual(t, "Should have parsed input.go correctly", expected, testutils.GetLastPrinted())
}

func Test_getFileFromArgs_shouldReturnFileName(t *testing.T) {
	defer teardown()
	setup()
	inputString := "input.go"
	addToArgs("-f=" + inputString)

	stringResult := getFileFromArgs()

	assertStringEqual(t, "Should have parsed -f flag correctly", inputString, stringResult)
}

func Test_getVariablesFromFunction_shouldParseTypedVariable(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
func main() {
	var typedVariable string
}`
	funStmt, _ := parseSource(source).Decls[0].(*ast.FuncDecl)
	// ast.Print(fset, funStmt.Body)

	findFunctionVariables(funStmt.Body)

	if len(functionVarsFound) != 1 {
		t.Fatalf("Did not find expected amount of vars, got %d ", len(functionVarsFound))
	}
	if exp := "typedVariable"; exp != functionVarsFound[0] {
		t.Errorf("Expected variable %q, got %q ", exp, functionVarsFound[0])
	}
}

func Test_getVariablesFromFunction_shouldReturnCompilerVariable(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
func main() {
	compilerTypedVar := "printed"
}`
	funStmt, _ := parseSource(source).Decls[0].(*ast.FuncDecl)

	findFunctionVariables(funStmt.Body)

	if len(functionVarsFound) != 1 {
		t.Fatalf("Did not find expected amount of vars, got %d ", len(functionVarsFound))
	}
	if exp := "compilerTypedVar"; exp != functionVarsFound[0] {
		t.Errorf("Expected variable %q, got %q ", exp, functionVarsFound[0])
	}
}

func Test_getVariablesFromFunction_shouldReturnParameter(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
func someFunc(parameter string) {
}`
	funStmt, _ := parseSource(source).Decls[0].(*ast.FuncDecl)
	// ast.Print(fset, funStmt)

	findParameters(funStmt.Type)

	if len(functionVarsFound) != 1 {
		t.Fatalf("Did not find expected amount of vars, got %d ", len(functionVarsFound))
	}
	if exp := "parameter"; exp != functionVarsFound[0] {
		t.Errorf("Expected variable %q, got %q ", exp, functionVarsFound[0])
	}
}

func Test_getVariablesFromFunction_shouldIgnoreSettingOfAField(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
var field string
func main() {
	v := "local"
	field = "member"
} `
	funStmt, _ := parseSource(source).Decls[1].(*ast.FuncDecl)

	findFunctionVariables(funStmt.Body)

	if len(functionVarsFound) != 1 {
		t.Fatalf("Did not find expected amount of vars, got %d ", len(functionVarsFound))
	}
}

func parseSource(src string) *ast.File {
	fset = token.NewFileSet() // positions are relative to fset
	f, _ := parser.ParseFile(fset, "", src, 0)
	return f
}

func Test_findFieldsFromFile_shouldFindTypedField(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
var variable string
`
	fileAst = parseSource(source)
	extractFieldsFromAst()

	if len(fieldsFound) != 1 {
		t.Fatalf("Did not find expected amount of vars, got %d ", len(fieldsFound))
	}
	if exp := "variable"; exp != fieldsFound[0] {
		t.Errorf("Expected variable %q, got %q ", exp, fieldsFound[0])
	}
}

func Test_findFieldsFromFile_shouldFindCompilerField(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
var fieldName = "string"
`
	fileAst = parseSource(source)
	extractFieldsFromAst()

	if len(fieldsFound) != 1 {
		t.Fatalf("Did not find expected amount of vars, got %d ", len(fieldsFound))
	}
	if exp := "fieldName"; exp != fieldsFound[0] {
		t.Errorf("Expected variable %q, got %q ", exp, fieldsFound[0])
	}
}

func Test_findFunctionsInFile_shouldFindSingleFunction(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
func funcName() {
	var variableName string
}`

	fileAst = parseSource(source)
	extractFunctionsFromAst()

	if len(functionsFound) != 1 {
		t.Fatalf("Did not find expected amount of functions, got %d ", len(functionsFound))
	}
	assertFunctionEquals(t, functionsFound[0], function{
		bodyStart: 2,
		bodyEnd:   4,
		variables: []string{"variableName"},
	})
}

func Test_findFunctionsInFile_shouldFindSingleFunctions(t *testing.T) {
	defer teardown()
	setup()
	source := `package main
func funcName() {
	var variableName string
}
func otherFunc() {
	otherVar := "a"
}`
	expectedFunctions := []function{
		function{
			bodyStart: 2,
			bodyEnd:   4,
			variables: []string{"variableName"},
		},
		function{
			bodyStart: 5,
			bodyEnd:   7,
			variables: []string{"otherVar"},
		},
	}

	fileAst = parseSource(source)
	extractFunctionsFromAst()

	for i, expFunc := range expectedFunctions {
		assertFunctionEquals(t, expFunc, functionsFound[i])
	}
}

func assertFunctionEquals(t *testing.T, m1, m2 function) {

	if m1.bodyStart != m2.bodyStart {
		t.Errorf("Expected function to start at line %d, was %d", m2.bodyStart, m1.bodyStart)
	}
	if m1.bodyEnd != m2.bodyEnd {
		t.Errorf("Expected function to end at line %d, was %d", m2.bodyEnd, m1.bodyEnd)
	}
	if len(m1.variables) != len(m2.variables) {
		t.Fatalf("Expected function to have %d vars, where %d.\n Vars: %q", len(m2.variables), len(m1.variables), m1.variables)
	}
	for j, cVar := range m1.variables {
		if cVar != m2.variables[j] {
			t.Errorf("Expected function to have var %q at [%d], but had %q", m2.variables[j], j, cVar)
		}
	}
}

// Helpers
//--------------------------------------------------------------------------------
func assertStringEqual(t *testing.T, message, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("%s\nExpected: %q\nGot     : %q", message, s1, s2)
	}
}

func assertIntEqual(t *testing.T, message string, s1, s2 int) {
	if s1 != s2 {
		t.Errorf("%s\nExpected: %q\nGot     : %q", message, s1, s2)
	}
}

func addToArgs(args ...string) {
	os.Args = append(os.Args, args...)
}
