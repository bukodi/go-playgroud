package fileio

import (
	"fmt"
	"testing"

	"github.com/Microsoft/go-winio"
)

func TestNamedPipe(t *testing.T) {
	path := `\\.\pipe\testPipe`
	listener, err := winio.ListenPipe(path, nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Hi server: %v\n", listener)
}
