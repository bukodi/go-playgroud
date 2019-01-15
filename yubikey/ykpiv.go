package main

import (
	"fmt"

	"pault.ag/go/ykpiv"
)

func main() {
	yubikey, err := ykpiv.New(ykpiv.Options{
		// Verbose: true,
		Reader: "Yubico Yubikey NEO U2F+CCID 01 00",
	})
	if err != nil {
		panic(err)
	}
	defer yubikey.Close()

	version, err := yubikey.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Application version %s found.\n", version)
}
