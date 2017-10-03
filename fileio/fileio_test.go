package fileio

import (
	_ "crypto/sha1"
	"fmt"
	_ "io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var filesBySHA1 map[string]os.FileInfo

func TestFileIO(t *testing.T) {

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

	fmt.Println("Hello")
}

func processDir(name string) {

}

func processFile(name string) {

}
