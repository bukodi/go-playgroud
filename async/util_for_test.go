package main

import (
	"strings"
	"testing"
)

type testCase struct {
	langs     []string
	greetings []string
	errs      []error
}

func createTestCase(name string, langs ...string) testCase {
	testCase := testCase{
		langs:     langs,
		greetings: make([]string, 0),
		errs:      make([]error, 0),
	}
	for _, lang := range langs {
		greeting, err := generateHello(name, lang)
		if err != nil {
			testCase.errs = append(testCase.errs, err)
		} else if greeting != "" {
			testCase.greetings = append(testCase.greetings, greeting)
		}
	}

	return testCase
}

func (tc testCase) checkResults(t *testing.T, greetings []string, errs []error) {
	expectedMsgs := messagesToString(tc.greetings, false)
	actualMsgs := messagesToString(greetings, false)
	if expectedMsgs != actualMsgs {
		t.Errorf("Greetings got: %v, want: %v", actualMsgs, tc.greetings)
	}
	expectedErrs := errorsToString(tc.errs, false)
	actualErrs := errorsToString(errs, false)
	if expectedErrs != actualErrs {
		t.Errorf("Errors got: %v, want: %v", actualErrs, expectedErrs)
	}
}

func errorsToString(errs []error, skipNils bool) string {
	parts := make([]string, 0)
	for _, err := range errs {
		if err != nil {
			parts = append(parts, err.Error())
		} else if !skipNils {
			parts = append(parts, "nil")
		} else {
			//Skip nil values
		}
	}
	return strings.Join(parts, ",")
}
func messagesToString(messages []string, skipEmpty bool) string {
	parts := make([]string, 0)
	for _, msg := range messages {
		if msg != "" {
			parts = append(parts, msg)
		} else if !skipEmpty {
			parts = append(parts, "")
		} else {
			//Skip nil values
		}
	}
	return strings.Join(parts, ",")
}
func errEquals(err1, err2 error) bool {
	if err1 == nil && err2 == nil {
		return true
	}
	if err1 != nil && err2 != nil && err1.Error() == err2.Error() {
		return true
	}
	return false
}
