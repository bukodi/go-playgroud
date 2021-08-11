package main

import (
	"time"
)

type Greeter interface {
	// This is the simples use case
	SayHello(name string) (greeting string)

	SayLocalHello(name string, lang string) (greeting string, err error)

	// Response is a stream
	SayMultiLangHello(name string, langs ...string) (greetings []string, errs []error)
}

type SyncGreeterImpl struct {
	delay time.Duration
}

var _ Greeter = &SyncGreeterImpl{}

func (g SyncGreeterImpl) SayHello(name string) (greeting string) {
	message, _ := g.SayLocalHello(name, "en")
	return message
}

func (g SyncGreeterImpl) SayLocalHello(name string, lang string) (greeting string, err error) {
	time.Sleep(g.delay)
	return generateLocalHello(name, lang)
}

func (g SyncGreeterImpl) SayMultiLangHello(name string, langs ...string) (greetings []string, errs []error) {
	greetings = make([]string, 0)
	errs = make([]error, 0)
	for _, lang := range langs {
		msg, err := g.SayLocalHello(name, lang)
		if err != nil {
			errs = append(errs, err)
		} else {
			greetings = append(greetings, msg)
		}
	}
	return greetings, errs
}
