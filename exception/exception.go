package main

import (
	"fmt"
)

func main() {
	ret, err := safeDiv2(10, 0)
	fmt.Printf("a / b = %d (Error: %v)", ret, err)
}

func safeDiv(a int, b int) (ret int, err error) {
	defer func() {
		if e := recover(); e != nil {
			e1, isErr := e.(error)
			if isErr {
				err = e1
			} else {
				err = fmt.Errorf("Div failed: %v", e)
			}
		}
	}()
	return div(a, b), nil
}

func safeDiv2(a int, b int) (ret int, err error) {
	defer func() { err = ToError(recover()) }()
	return div(a, b), nil
}

func div(a int, b int) int {
	return a / b
}

func ToError(recovered interface{}) error {
	if recovered == nil {
		return nil
	}
	err, isErr := recovered.(error)
	if isErr {
		return err
	} else {
		return fmt.Errorf("%v", recovered)
	}

}
