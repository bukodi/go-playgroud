package fileio

import (
	"flag"
	"fmt"
	os "os"
	"path/filepath"
)

func Main(){
	fmt.Println( os.Args)
	var (
		dbpath *string = flag.String("db", "./pathlist", "Path to cfg file")
	)
	flag.Parse()

	// f, err := os.Open(".");

	fmt.Println( filepath.Abs("."))
	fmt.Println("Hello")
	fmt.Println(*dbpath)
}
