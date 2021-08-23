package main

import (
	"context"
)

func WrapToAsyncGreeter(ctx context.Context, syncGreeter Greeter) AsyncGreeter {
	wrapper := syncToAsync{
		ctx:  ctx,
		impl: syncGreeter,
	}
	return &wrapper
}

var _ AsyncGreeter = &syncToAsync{}

type syncToAsync struct {
	impl Greeter
	ctx  context.Context
}

func (w syncToAsync) SayHello(ctx context.Context, name string) <-chan string {
	internalCh := make(chan string)
	go func() {
		defer close(internalCh)
		internalCh <- w.impl.SayHello(name)
	}()
	greetingCh := make(chan string)
	go func() {
		defer close(greetingCh)
		select {
		case <-ctx.Done():
			return
		case msg := <-internalCh:
			greetingCh <- msg
		}
	}()
	return greetingCh
}

func (w syncToAsync) SayLocaleHello(ctx context.Context, name string, lang string) (<-chan string, <-chan error) {
	internalMsgCh := make(chan string)
	internalErrCh := make(chan error)
	go func() {
		defer close(internalMsgCh)
		defer close(internalErrCh)
		msg, err := w.impl.SayLocaleHello(name, lang)
		if err != nil {
			internalErrCh <- err
		} else {
			internalMsgCh <- msg
		}
	}()

	greetingCh := make(chan string)
	errCh := make(chan error)
	go func() {
		defer close(greetingCh)
		defer close(errCh)

		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		case msg, ok := <-internalMsgCh:
			if ok {
				greetingCh <- msg
			} else {
				return
			}
		case err, ok := <-internalErrCh:
			if ok {
				errCh <- err
			} else {
				return
			}
		}
	}()
	return greetingCh, errCh
}

func (w syncToAsync) SayMultiLangHello(ctx context.Context, name string, langCh <-chan string) (<-chan string, <-chan error) {
	internalMsgCh := make(chan string)
	internalErrCh := make(chan error)
	go func() {
		defer close(internalMsgCh)
		defer close(internalErrCh)

		langs := make([]string, 0)
		for lang := range langCh {
			langs = append(langs, lang)
		}
		msgs, errs := w.impl.SayMultiLangHello(name, langs...)
		for _, err := range errs {
			if err != nil {
				internalErrCh <- err
			}
		}
		for _, msg := range msgs {
			internalMsgCh <- msg
		}
	}()

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
			case msg, ok := <-internalMsgCh:
				if ok {
					greetingCh <- msg
				} else {
					break loop
				}
			case err, ok := <-internalErrCh:
				if ok {
					errCh <- err
				} else {
					break loop
				}
			}
		}
	}()
	return greetingCh, errCh
}
