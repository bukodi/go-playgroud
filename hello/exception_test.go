package main

import (
	"errors"
	"fmt"
	"testing"
)

func fnHighLevel() error {
	_ = fnLowLevel()
	return nil
}

func fnLowLevel() error {
	panic(errors.New("Ezt kapd el"))
}

func TestException(t *testing.T) {
	err := fnHighLevel()
	fmt.Printf("%#v", err)
}
