package Mylibrary

import "fmt"

// Add adds two integers and returns the result.
func Add(a, b int) int {
	return a + b
}

// Multiply multiplies two integers and returns the result.
func Multiply(a, b int) int {
	return a * b
}

var name string
func Rohit(s string)  {
	if name == "" {
		panic("Plz set name")
	}
	fmt.Println(name," Says",s)	
}

func SetName(s string){
	name= s 
}