package main

import (
	"context"
	"time"
)

type AsyncGreeterImpl struct {
	delay time.Duration
}

var _ AsyncGreeter = &AsyncGreeterImpl{}

func (a AsyncGreeterImpl) SayHello(ctx context.Context, name string) <-chan string {
	greetingCh := make(chan string)
	go func() {
		defer close(greetingCh)
		select {
		case <-ctx.Done():
			return
		case <-time.After(a.delay):
		}
		msg, _ := generateHello(name, "en")
		greetingCh <- msg
	}()
	return greetingCh
}

func (a AsyncGreeterImpl) SayLocaleHello(ctx context.Context, name string, lang string) (<-chan string, <-chan error) {
	greetingCh := make(chan string)
	errCh := make(chan error)
	go func() {
		defer close(greetingCh)
		defer close(errCh)

		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		case <-time.After(a.delay):
		}

		msg, err := generateHello(name, lang)
		if err != nil {
			errCh <- err
		} else {
			greetingCh <- msg
		}
	}()
	return greetingCh, errCh
}

func (a AsyncGreeterImpl) SayMultiLangHello(ctx context.Context, name string, langCh <-chan string) (greetingsCh <-chan string, err <-chan error) {
	greetingCh := make(chan string)
	errCh := make(chan error)
	go func() {
		defer close(greetingCh)
		defer close(errCh)
	loop:
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				break loop
			case lang, ok := <-langCh:
				if !ok {
					break loop
				}
				msg, err := generateHello(name, lang)
				if err != nil {
					errCh <- err
				} else {
					greetingCh <- msg
				}
			}
		}
	}()
	return greetingCh, errCh
}
