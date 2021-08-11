package main

import (
	"fmt"
	"sync"
	"testing"
)

func worker(nameCh <-chan string) (greetingsCh <-chan string) {
	greetingCh := make(chan string)
	go func() {
		defer close(greetingCh)
	loop:
		for {
			select {
			case name, ok := <-nameCh:
				if !ok {
					break loop
				}
				greetingCh <- fmt.Sprintf("Hello %s!", name)
			}
		}
	}()
	return greetingCh
}

func TestDeadlock(t *testing.T) {
	var wg sync.WaitGroup
	nameCh := make(chan string)
	msgCh := worker(nameCh)
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
				t.Logf("received : %s", msg)
			}
		}
	}()

	for _, name := range []string{"Alice", "Bob", "Carole"} {
		t.Logf("sending  : %s", name)
		nameCh <- name
	}

	// WARNING: Is you skip this line, receive this:
	// fatal error: all goroutines are asleep - deadlock!
	close(nameCh)

	wg.Wait()
	t.Log("Completed")

}
