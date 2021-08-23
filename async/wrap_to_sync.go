package main

import (
	"context"
	"sync"
)

func WrapToSyncGreeter(ctx context.Context, asyncGreeter AsyncGreeter) Greeter {
	wrapper := asyncToSync{
		ctx:  ctx,
		impl: asyncGreeter,
	}
	return &wrapper
}

type asyncToSync struct {
	impl AsyncGreeter
	ctx  context.Context
}

var _ Greeter = &asyncToSync{}

func (w asyncToSync) SayHello(name string) (greeting string) {
	return <-w.impl.SayHello(w.ctx, name)
}

func (w asyncToSync) SayLocaleHello(name string, lang string) (greeting string, err error) {
	msgCh, errCh := w.impl.SayLocaleHello(w.ctx, name, lang)
	select {
	case <-w.ctx.Done():
		return "", w.ctx.Err()
	case msg := <-msgCh:
		return msg, nil
	case err := <-errCh:
		return "", err
	}
}

func (w asyncToSync) SayMultiLangHello(name string, langs ...string) (greetings []string, errs []error) {
	langCh := make(chan string, len(langs))
	msgCh, errCh := w.impl.SayMultiLangHello(w.ctx, name, langCh)

	var wg sync.WaitGroup
	greetings = make([]string, 0)
	errs = make([]error, 0)

	wg.Add(1)
	go func() {
		defer wg.Done()
	loop:
		for {
			select {
			case <-w.ctx.Done():
				errs = append(errs, w.ctx.Err())
				break loop
			case msg, ok := <-msgCh:
				if !ok {
					break loop
				}
				greetings = append(greetings, msg)
			case err, ok := <-errCh:
				if !ok {
					break loop
				}
				errs = append(errs, err)
			}
		}
	}()

	for _, lang := range langs {
		langCh <- lang
	}
	close(langCh)

	wg.Wait()
	return greetings, errs
}
