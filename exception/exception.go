package main

import (
	"fmt"
)

func main2() {
	ret, err := safeDiv2(10, 1)
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

func safeDiv3(a int, b int) (ret int, err error) {
	defer func() {
		err = ToError(recover())
		if err != nil {
			fmt.Printf("Catch error: %v\n", err)
		}
		fmt.Printf("Finally\n")
	}()
	return div(a, b), nil
}

func main() {
	a := 10
	b := 0
	ret1, err1 := func() (ret int, err error) {
		defer func() {
			err = ToError(recover())
			if err != nil {
				fmt.Printf("Catch error: %v\n", err)
			}
			fmt.Printf("Finally\n")
		}()

		return a / b, nil
	}()
	fmt.Printf("%d %v\n", ret1, err1)
}

func divWrapper(a int, b int) int {
	return div(a, b)
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
