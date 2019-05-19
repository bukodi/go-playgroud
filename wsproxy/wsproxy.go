package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"

	"github.com/koding/websocketproxy"
)

var (
	flagBackend = flag.String("backend", "", "Backend URL for proxying")
)

func main() {
	u, err := url.Parse(*flagBackend)
	u, err = url.Parse("wss://echo.websocket.org")
	if err != nil {
		log.Fatalln(err)
	}

	err = http.ListenAndServe(":8888", websocketproxy.NewProxy(u))
	if err != nil {
		log.Fatalln(err)
	}
}
