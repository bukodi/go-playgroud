package suffixtree

// https://eli.thegreenplace.net/2016/suffix-arrays-in-the-go-standard-library/

import (
	"fmt"
	"index/suffixarray"
	"strings"
	"testing"
)

func TestTree(t *testing.T) {
	t.Log("Hello")

	words := []string{
		"banana",
		"apple",
		"pear",
		"tangerine",
		"orange",
		"lemon",
		"peach",
		"persimmon",
	}

	data := []byte("\x00" + strings.Join(words, "\x00") + "\x00")
	sa := suffixarray.New(data)

	indices := sa.Lookup([]byte("an"), -1)

	// Reconstruct matches from indices found by Lookup.
	for _, idx := range indices {
		fmt.Println(getStringFromIndex(data, idx))
	}
}

func getStringFromIndex(data []byte, index int) string {
	var start, end int
	for i := index - 1; i >= 0; i-- {
		if data[i] == 0 {
			start = i + 1
			break
		}
	}
	for i := index + 1; i < len(data); i++ {
		if data[i] == 0 {
			end = i
			break
		}
	}
	return string(data[start:end])
}
