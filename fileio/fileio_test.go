package fileio

import (
	_ "crypto/sha1"
	_ "fmt"
	_ "io/ioutil"
	"os"
	_ "path/filepath"
	"testing"
	"crypto/x509"
	"fmt"
)

var filesBySHA1 map[string]os.FileInfo

func TestFileIO(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmdName", "-db", "./cica"}

	certPool, err := x509.SystemCertPool()
	fmt.Println(err)
	fmt.Println(certPool)

	Main()

	/*
	list := make([]string, 5)
	err := filepath.Walk("/opt/google/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		list = append(list, path)
		return nil
	})
	if err != nil {
		fmt.Printf("walk error [%v]\n", err)
	}

	for _, f := range list {
		fmt.Println(f)
	}

	fmt.Println("Hello") */
}

func processDir(name string) {

}

func processFile(name string) {

}
