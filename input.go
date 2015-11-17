package main

import "fmt"

var someFieldToPrint = "value of someFieldToPrint"
var someOtherField = "value of someOtherField"

func someFunc() {
	stringToPrint := "printed"
	a, b := multiFunc()
	fmt.Println(stringToPrint)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(someFieldToPrint)
	fmt.Println(someOtherField)
}

func multiFunc() (string, string) {
	a := "a"
	return a, "b"
}
