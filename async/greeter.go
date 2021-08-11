package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func main() {
	syncImpl := SyncGreeterImpl{delay: time.Millisecond * 400}
	fmt.Println(syncImpl.SayHello("Alice"))
	fmt.Println(syncImpl.SayLocalHello("Alice", "hu"))
	fmt.Println(syncImpl.SayLocalHello("Alice", "xx"))
	fmt.Println(syncImpl.SayMultiLangHello("Alice", "en", "xx", "fr", "yy", "es"))

	asyncImpl := AsyncGreeterImpl{delay: time.Millisecond * 400}
	fmt.Println(<-asyncImpl.SayHello(context.Background(), "Alice"))

	//fmt.Println( <-asyncImpl.SayLocalHello(context.Background(), "Alice", "hu"))
}

func generateLocalHello(name string, lang string) (greeting string, err error) {
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
