package main

import "context"

type Greeter interface {
	// This is the simples use case
	SayHello(name string) (greeting string)

	SayLocaleHello(name string, lang string) (greeting string, err error)

	// Response is a stream
	SayMultiLangHello(name string, langs ...string) (greetings []string, errs []error)
}

type AsyncGreeter interface {
	// This is the simplest use case
	SayHello(ctx context.Context, name string) (greetingCh <-chan string)

	// Response with optional error
	SayLocaleHello(ctx context.Context, name string, lang string) (greetingCh <-chan string, errCh <-chan error)

	// Input and output is streams
	SayMultiLangHello(ctx context.Context, name string, langCh <-chan string) (greetingsCh <-chan string, err <-chan error)
}
