package main

import (
	"context"
	"time"
)

type AsyncGreeter interface {
	// This is the simples use case
	SayHello(ctx context.Context, name string) (greetingCh <-chan string)

	// This is the simples use case
	SayLocalHello(ctx context.Context, name string, lang string) (greetingCh <-chan string, errCh <-chan error)

	// Response is a stream
	SayMultiLangHello(ctx context.Context, name string, lang ...string) (greetingsCh <-chan string, err <-chan error)
}

type AsyncGreeterImpl struct {
	delay time.Duration
}

func (a AsyncGreeterImpl) SayHello(ctx context.Context, name string) <-chan string {
	greetingCh := make(chan string)
	go func() {
		select {
		case <-ctx.Done():
			close(greetingCh)
			return
		case <-time.After(a.delay):
		}
		msg, _ := generateLocalHello(name, "en")
		greetingCh <- msg
		close(greetingCh)
	}()
	return greetingCh
}

func (a AsyncGreeterImpl) SayLocalHello(ctx context.Context, name string, lang string) (<-chan string, <-chan error) {
	greetingCh := make(chan string)
	errCh := make(chan error)
	go func() {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			close(greetingCh)
			close(errCh)
			return
		case <-time.After(a.delay):
		}
		msg, err := generateLocalHello(name, lang)
		if err != nil {
			errCh <- err
		} else {
			greetingCh <- msg
		}
		close(greetingCh)
		close(errCh)
	}()
	return greetingCh, errCh
}

func (a AsyncGreeterImpl) SayMultiLangHello(ctx context.Context, name string, lang ...string) (greetingsCh <-chan string, err <-chan error) {
	panic("implement me")
}

var _ AsyncGreeter = &AsyncGreeterImpl{}
