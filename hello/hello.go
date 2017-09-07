package main

import (
	"fmt"
	"reflect"
	_ "reflect"
)

type Cica struct {
	Fej   string
	Farok string
}

var cica Cica

func main() {
	fmt.Printf("Hello, world.\n")
	t := reflect.TypeOf(cica)
	fmt.Println(t)
}
