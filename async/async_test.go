package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

var asyncImpl AsyncGreeter = AsyncGreeterImpl{delay: time.Millisecond * 100}

func TestAsyncSayHello(t *testing.T) {
	t.Run("AsyncGreeterImpl", func(t *testing.T) {
		testAsyncSayHello(t, asyncImpl)
	})
	wrappedSync := WrapToAsyncGreeter(context.Background(), SyncGreeterImpl{delay: time.Millisecond * 100})
	t.Run("WrapToAsyncGreeter", func(t *testing.T) {
		testAsyncSayHello(t, wrappedSync)
	})
}
func testAsyncSayHello(t *testing.T, impl AsyncGreeter) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	expectedMsg, _ := generateHello("Alice", "en")
	msg := <-impl.SayHello(ctx, "Alice")
	if expectedMsg != msg {
		t.Errorf("actual: %v, expected: %v", msg, expectedMsg)
	}
}

func TestAsyncSayHelloWithTimeout(t *testing.T) {
	t.Run("AsyncGreeterImpl", func(t *testing.T) {
		testAsyncSayHelloWithTimeout(t, asyncImpl)
	})
	wrappedSync := WrapToAsyncGreeter(context.Background(), SyncGreeterImpl{delay: time.Millisecond * 100})
	t.Run("WrapToAsyncGreeter", func(t *testing.T) {
		testAsyncSayHelloWithTimeout(t, wrappedSync)
	})
}

func testAsyncSayHelloWithTimeout(t *testing.T, impl AsyncGreeter) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()
	expectedMsg := ""
	msg := <-impl.SayHello(ctx, "Alice")
	if expectedMsg != msg {
		t.Errorf("actual: %v, expected: %v", msg, expectedMsg)
	}
}
func TestAsyncSayMultiLangHello(t *testing.T) {
	t.Run("AsyncGreeterImpl", func(t *testing.T) {
		testAsyncSayMultiLangHello(t, asyncImpl)
	})
	wrappedSync := WrapToAsyncGreeter(context.Background(), SyncGreeterImpl{delay: time.Millisecond * 100})
	t.Run("WrapToAsyncGreeter", func(t *testing.T) {
		testAsyncSayMultiLangHello(t, wrappedSync)
	})
}

func testAsyncSayMultiLangHello(t *testing.T, impl AsyncGreeter) {
	var wg sync.WaitGroup
	langCh := make(chan string)
	msgCh, errCh := impl.SayMultiLangHello(context.Background(), "Alice", langCh)
	wg.Add(1)
	go func() {
		defer wg.Done()
	loop:
		for {
			select {
			case msg, ok := <-msgCh:
				if !ok {
					break loop
				}
				t.Log("Received: ", msg, nil)
			case err, ok := <-errCh:
				if !ok {
					break loop
				}
				t.Log("Received: ", nil, err)
			}
		}
	}()

	for _, lang := range []string{"en", "xx", "fr", "yy", "es"} {
		time.Sleep(time.Millisecond * 100)
		t.Log("Sending: ", lang)
		langCh <- lang
	}
	close(langCh)

	wg.Wait()
	t.Log("Completed")

}
