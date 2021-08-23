package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func generateHello(name string, lang string) (greeting string, err error) {
	switch lang {
	case "en":
		return fmt.Sprintf("Hello %s!", name), nil
	case "fr":
		return fmt.Sprintf("Bonjour %s!", name), nil
	case "es":
		return fmt.Sprintf("Hola %s", name), nil
	case "hu":
		return fmt.Sprintf("Szia %s!", name), nil
	default:
		return "", errors.Errorf("unsupported lang: %s", lang)
	}
}
