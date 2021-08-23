package main

import (
	"context"
	"strings"
	"testing"
	"time"
)

var syncImpl = SyncGreeterImpl{delay: time.Millisecond * 100}
var wrappedAsync = WrapToSyncGreeter(context.Background(), AsyncGreeterImpl{delay: time.Millisecond * 100})

func TestSyncSayHello(t *testing.T) {
	t.Run("SyncGreeterImpl", func(t *testing.T) {
		testSyncSayHello(t, syncImpl)
	})
	t.Run("WrapToSyncGreeter", func(t *testing.T) {
		testSyncSayHello(t, wrappedAsync)
	})
}
func testSyncSayHello(t *testing.T, impl Greeter) {
	expected, _ := generateHello("Alice", "en")
	actual := impl.SayHello("Alice")
	if expected != actual {
		t.Errorf("SayHello=%v, expected: %v", actual, expected)
	}
}

func TestSyncSayLocaleHello(t *testing.T) {
	t.Run("SyncGreeterImpl", func(t *testing.T) {
		testSyncSayLocaleHello(t, syncImpl)
	})
	t.Run("WrapToSyncGreeter", func(t *testing.T) {
		testSyncSayLocaleHello(t, wrappedAsync)
	})
}

func testSyncSayLocaleHello(t *testing.T, impl Greeter) {
	langs := []string{"en", "de", "xx", "fr", "hu", "yy"}
	for _, lang := range langs {
		t.Run(lang, func(t *testing.T) {
			t.Parallel()
			expectedMsg, expectedErr := generateHello("Alice", lang)
			gotGreeting, err := impl.SayLocaleHello("Alice", lang)
			if !errEquals(err, expectedErr) {
				t.Errorf("SayLocaleHello() error = %v, wantErr %v", err, expectedErr)
				return
			}
			if gotGreeting != expectedMsg {
				t.Errorf("SayLocaleHello() gotGreeting = %v, want %v", gotGreeting, expectedErr)
			}
		})
	}
}

func TestSyncSayMultiHello(t *testing.T) {
	t.Run("SyncGreeterImpl", func(t *testing.T) {
		testSyncSayMultiHello(t, syncImpl)
	})
	t.Run("WrapToSyncGreeter", func(t *testing.T) {
		testSyncSayMultiHello(t, wrappedAsync)
	})
}

func testSyncSayMultiHello(t *testing.T, impl Greeter) {
	tests := []testCase{
		createTestCase("Alice", "hu", "en", "xx"),
		createTestCase("Alice"),
		createTestCase("Alice", "yy", "es", "fr"),
	}
	for _, test := range tests {
		t.Run(strings.Join(test.langs, ","), func(t *testing.T) {
			t.Parallel()
			gotGreetings, errs := impl.SayMultiLangHello("Alice", test.langs...)
			test.checkResults(t, gotGreetings, errs)
		})
	}
}
