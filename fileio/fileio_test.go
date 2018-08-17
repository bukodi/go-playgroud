package fileio

import (
	_ "crypto/sha1"
	"crypto/x509"
	"fmt"
	_ "io/ioutil"
	"os"
	"path/filepath"
	"testing"
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

	list := make([]string, 5)
	err = filepath.Walk("/opt/", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Println("Dir:" + path)
			return nil
		}
		list = append(list, path)
		fmt.Println(path)
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

func processDirs(oldBase string, newBase string) {
	filepath.Walk(newBase, func(newPath string, newInfo os.FileInfo, err error) error {
		rel, _ := filepath.Rel(newBase, newPath)
		oldPath := filepath.Join(oldBase, rel)
		oldInfo, err := os.Stat(oldPath)

		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", oldPath, err)
			return nil
		}
		err = os.Chmod(newPath, func() os.FileMode {
			if newInfo.IsDir() {
				return 0777
			} else {
				return 0666
			}
		}())
		if err != nil {
			fmt.Printf("chmod failed (%q) : %v\n", newPath, err)
			return nil
		}
		err = os.Chown(newPath, os.Getuid(), os.Getuid())
		if err != nil {
			fmt.Printf("chown failed (%q) : %v\n", newPath, err)
			return nil
		}

		err = os.Chtimes(newPath, oldInfo.ModTime(), oldInfo.ModTime())
		if err != nil {
			fmt.Printf("chtimes failed (%q) : %v\n", newPath, err)
			return nil
		}

		//fmt.Printf("%q - orig date: %+v\n", newPath, oldInfo.ModTime().String())
		return nil
	})
}

func TestProcessDirs(t *testing.T) {
	processDirs("/home/bukodi/Arduino", "/tmp/Arduino")
}

func TestWalk(t *testing.T) {
	dir := "/home/bukodi/Arduino"
	subDirToSkip := "skip" // dir/to/walk/skip

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", dir, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited cic file: %q\n", path)
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", dir, err)
	}
}
