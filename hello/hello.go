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

	var ia5 *[5]int
	ia5 = new([5]int)
	fmt.Print(ia5)

	var s1 struct {
		v1 int
		v2 int
	}
	s1.v1 = 5
	s1.v2 = 6
	fmt.Print(s1)

	ca, _ := Asset("cica.txt")
	fmt.Println(ca)
	c, _ := cicaTxt()
	fmt.Println(c)
}
