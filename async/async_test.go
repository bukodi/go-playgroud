package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestAsyncGreeterImpl_SayMultiLangHello(t *testing.T) {
	asyncImpl := AsyncGreeterImpl{delay: time.Millisecond * 400}

	var wg sync.WaitGroup
	langCh := make(chan string)
	msgCh, errCh := asyncImpl.SayMultiLangHello(context.Background(), "Alice", langCh)
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
