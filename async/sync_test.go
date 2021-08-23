package main

import (
	"context"
	"runtime"
	"strings"
	"testing"
	"time"
)

var syncImpl = SyncGreeterImpl{delay: time.Millisecond * 100}
var wrappedAsync = WrapToSyncGreeter(context.Background(), AsyncGreeterImpl{delay: time.Millisecond * 100})

func TestSyncSayHello(t *testing.T)         { testSyncSayHello(t, syncImpl) }
func TestWrappedAsyncSayHello(t *testing.T) { testSyncSayHello(t, wrappedAsync) }

func testSyncSayHello(t *testing.T, impl Greeter) {
	defer checkGoroutineLeakage(t, runtime.NumGoroutine())

	expected, _ := generateHello("Alice", "en")
	actual := impl.SayHello("Alice")
	if expected != actual {
		t.Errorf("SayHello=%v, expected: %v", actual, expected)
	}
}

func TestSyncSayLocaleHello(t *testing.T)         { testSyncSayLocaleHello(t, syncImpl) }
func TestWrappedAsyncSayLocaleHello(t *testing.T) { testSyncSayLocaleHello(t, wrappedAsync) }

func testSyncSayLocaleHello(t *testing.T, impl Greeter) {
	defer checkGoroutineLeakage(t, runtime.NumGoroutine())

	langs := []string{"en", "de", "xx", "fr", "hu", "yy"}
	for _, lang := range langs {
		t.Run(lang, func(t *testing.T) {
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

func TestSyncSayMultiHello(t *testing.T)         { testSyncSayMultiHello(t, syncImpl) }
func TestWrappedAsyncSayMultiHello(t *testing.T) { testSyncSayMultiHello(t, wrappedAsync) }

func testSyncSayMultiHello(t *testing.T, impl Greeter) {
	defer checkGoroutineLeakage(t, runtime.NumGoroutine())

	base := runtime.NumGoroutine()
	tests := []testCase{
		createTestCase("Alice", "hu", "en", "xx"),
		createTestCase("Alice"),
		createTestCase("Alice", "yy", "es", "fr"),
	}
	for _, test := range tests {
		t.Run(strings.Join(test.langs, ","), func(t *testing.T) {
			//t.Parallel()
			gotGreetings, errs := impl.SayMultiLangHello("Alice", test.langs...)
			test.checkResults(t, gotGreetings, errs)
		})
	}
	runtime.Gosched()
	after := runtime.NumGoroutine()
	if after > base {
		t.Fail()
	}
}
