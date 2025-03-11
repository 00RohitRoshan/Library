package Mylibrary

import (
	"fmt"
	"os"
)

var name string

func File(s string) {
	checkName()
	file, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	fmt.Fprintln(file,name, " Says", s)
}

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
