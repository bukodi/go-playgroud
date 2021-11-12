package errorsandlogs

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type InvalidInputError struct {
	input string
}

func (iie *InvalidInputError) Error() string {
	return fmt.Sprintf("input cause error: %s", iie.input)
}

func fnLevel1(input string) error {
	if strings.Contains(input, "1") && strings.Contains(input, "error") {
		return &InvalidInputError{input: input}
	}
	err := fnLevel2(input)
	if err != nil {
		return fmt.Errorf("fnLevel2 call failed: %w", err)
	}
	return nil
}

func fnLevel2(input string) error {
	if strings.Contains(input, "2") && strings.Contains(input, "error") {
		return &InvalidInputError{input: input}
	}
	err := fnLevel3(input)
	if err != nil {
		return fmt.Errorf("fnLevel3 call failed: %w", err)
	}
	return nil
}

func fnLevel3(input string) error {
	if strings.Contains(input, "3") && strings.Contains(input, "error") {
		return &InvalidInputError{input: input}
	}
	return nil
}

func TestWrapping(t *testing.T) {
	err := fnLevel1("error 3")
	t.Log(err)

	unwrapped1 := errors.Unwrap(err)
	t.Logf("unwrapped1: %v", unwrapped1)
	unwrapped2 := errors.Unwrap(unwrapped1)
	t.Logf("unwrapped2: %v", unwrapped2)
	unwrapped3 := errors.Unwrap(unwrapped2)
	t.Logf("unwrapped3: %v", unwrapped3)
}

func TestAs(t *testing.T) {
	err := fnLevel1("error 3")
	t.Log(err)

	iie := new(InvalidInputError)
	if errors.As(err, &iie) {
		t.Log(iie)
	}
}

func TestIs(t *testing.T) {
	err := fnLevel1("error 3")
	t.Log(err)

	iie := new(InvalidInputError)
	if errors.Is(err, iie) {
		t.Log("Ok, this is an InvalidInputError")
	} else {
		t.Log("This is not an InvalidInputError")
	}
}
