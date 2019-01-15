package main

import (
	"net/http"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	acmeClient := acme.Client{
	//DirectoryURL: "https://acme.api.letsencrypt.org/directory",
	//DirectoryURL: "https://acme-staging.api.letsencrypt.org/directory",
	//DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
	}
	certManager := autocert.Manager{
		Client:     &acmeClient,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("bukodi.ddns.net"), //Your domain here
		Cache:      autocert.DirCache("/tmp/certs2"),          //Folder for storing certificates
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	server := &http.Server{
		Addr:      ":10443",
		TLSConfig: certManager.TLSConfig(),
	}

	go http.ListenAndServe("https://localhost:10443", certManager.HTTPHandler(nil))
	//http.HandleFunc("/", helloHandler)
	server.ListenAndServeTLS("", "") //Key and cert are coming from Let's Encrypt
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}
