package main

import (
	"time"
)

type SyncGreeterImpl struct {
	delay time.Duration
}

var _ Greeter = &SyncGreeterImpl{}

func (g SyncGreeterImpl) SayHello(name string) (greeting string) {
	message, _ := g.SayLocaleHello(name, "en")
	return message
}

func (g SyncGreeterImpl) SayLocaleHello(name string, lang string) (greeting string, err error) {
	time.Sleep(g.delay)
	return generateHello(name, lang)
}

func (g SyncGreeterImpl) SayMultiLangHello(name string, langs ...string) (greetings []string, errs []error) {
	greetings = make([]string, 0)
	errs = make([]error, 0)
	for _, lang := range langs {
		msg, err := g.SayLocaleHello(name, lang)
		if err != nil {
			errs = append(errs, err)
		} else {
			greetings = append(greetings, msg)
		}
	}
	return greetings, errs
}
