package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var fmtPrintf = fmt.Printf
var fileNamePtr *string
var fset *token.FileSet
var fileAst *ast.File

var fieldsFound []string
var functionsFound []function
var functionVarsFound []string

func init() {
	fileNamePtr = flag.String("f", "", "Filename to parse")
}

func main() {
	file := getFileFromArgs()

	parseFile(file)
	extractWantedFromAst()
	printFound()
}

func getFileFromArgs() string {
	flag.Parse()
	return *fileNamePtr
}

func parseFile(file string) {
	fset = token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		panic("error parsing file \"" + file + "\".")
	}
	fileAst = f
}

func extractWantedFromAst() {
	extractFieldsFromAst()
	extractFunctionsFromAst()
}

func printFound() {
	var groups []string

	fields := fieldsFound
	groups = append(groups, strings.Join(fields, " "))
	functions := functionsFound
	for _, currentFunction := range functions {
		varString := strings.Join(currentFunction.variables, " ")
		functionString := fmt.Sprintf("%d,%d,%s", currentFunction.bodyStart, currentFunction.bodyEnd, varString)

		groups = append(groups, functionString)
	}

	fmtPrintf(strings.Join(groups, "|"))
}

// FIELDS
//================================================================================

func extractFieldsFromAst() {
	for _, declaration := range fileAst.Decls {
		if genDecl, isGenDecl := declaration.(*ast.GenDecl); isGenDecl {
			for _, s := range genDecl.Specs {
				if value, isValueSpec := s.(*ast.ValueSpec); isValueSpec {
					addFoundField(value.Names[0].String())
				}
			}
		}
	}
}

func addFoundField(f string) {
	fieldsFound = append(fieldsFound, f)
}

// FUNCTIONS
//================================================================================

type function struct {
	bodyStart int      //Line the function starts at
	bodyEnd   int      //Line the function ends at
	variables []string //a list of variables
}

func extractFunctionsFromAst() {
	for _, declaration := range fileAst.Decls {
		if funcDecl, isFuncDecl := declaration.(*ast.FuncDecl); isFuncDecl {
			processFunction(funcDecl)
		}
	}
}

func processFunction(funcDecl *ast.FuncDecl) {
	m := function{}

	m.bodyStart = fset.Position(funcDecl.Pos()).Line
	m.bodyEnd = fset.Position(funcDecl.End()).Line
	m.variables = getFunctionVariables(funcDecl)

	addFoundFunctions(m)
}

func addFoundFunctions(m ...function) {
	functionsFound = append(functionsFound, m...)
}

// Variables
//--------------------------------------------------------------------------------

func getFunctionVariables(funcDecl *ast.FuncDecl) []string {
	functionVarsFound = nil
	findParameters(funcDecl.Type)
	findFunctionVariables(funcDecl.Body)
	return functionVarsFound
}

func findParameters(funcType *ast.FuncType) {
	for _, param := range funcType.Params.List {
		for _, name := range param.Names {
			addFoundFunctionVar(name.Name)
		}
	}
}

func findFunctionVariables(functionBlock *ast.BlockStmt) {
	for _, functionStmt := range functionBlock.List {
		if decl, isDecl := functionStmt.(*ast.DeclStmt); isDecl {
			getVariablesFromDeclStmt(decl)
		} else if assign, isAssign := functionStmt.(*ast.AssignStmt); isAssign {
			if assign.Tok != token.ASSIGN {
				getVariablesFromAssignStmt(assign)
			}
		}
	}
}

func getVariablesFromAssignStmt(assignStmt *ast.AssignStmt) {
	for _, itemLeftOfAssignment := range assignStmt.Lhs {
		if variable, isVarIdent := itemLeftOfAssignment.(*ast.Ident); isVarIdent {
			addFoundFunctionVar(variable.Name)
		}
	}
}

func getVariablesFromDeclStmt(decl *ast.DeclStmt) {
	if genDecl, isGenDecl := decl.Decl.(*ast.GenDecl); isGenDecl {
		for _, s := range genDecl.Specs {
			if value, isValueSpec := s.(*ast.ValueSpec); isValueSpec {
				for _, ident := range value.Names {
					addFoundFunctionVar(ident.Name)
				}
			}
		}
	}
}

func addFoundFunctionVar(v ...string) {
	functionVarsFound = append(functionVarsFound, v...)
}
