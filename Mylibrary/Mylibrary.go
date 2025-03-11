package Mylibrary

import "fmt"

// Add adds two integers and returns the result.
func Add(a, b int) int {
	checkName()
	return a + b
}

// Multiply multiplies two integers and returns the result.
func Multiply(a, b int) int {
	checkName()
	return a * b
}

var name string

func Console(s string) {
	checkName()
	fmt.Println(name, " Says", s)
}

func checkName() {
	if name == "" {
		panic("Plz set name")
	}
}

func SetName(s string) {
	name = s
}
