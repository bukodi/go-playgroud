package main

import (
	"fmt"
	_ "reflect"
	"reflect"
)

type Cica struct {
	Fej string
	Farok string
}

var cica Cica

func main() {
	fmt.Printf("Hello, world.\n")
	t := reflect.TypeOf(cica);
	fmt.Println(t)
	t := reflect.TypeOf(cica);
	fmt.Println(t)
}
