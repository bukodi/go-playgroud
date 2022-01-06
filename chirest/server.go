package chirest

import (
	"log"
	"net/http"
)

func StartServer(addr string, handler http.Handler) (*http.Server, error) {

	srv := &http.Server{Addr: addr, Handler: handler}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println("HTTP Server Error - ", err)
		}
	}()

	return srv, nil
}
