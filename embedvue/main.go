package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"sync"
)

//go:embed _ui/dist
var uiDist embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(uiDist, "_ui/dist")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fsys)
}

type countHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	fmt.Fprintf(w, "request URI %s\n", r.RequestURI)
	fmt.Fprintf(w, "request URL path  %s\n", r.URL.Path)
	fmt.Fprintf(w, "count is %d\n", h.n)
}

func main() {
	fs := http.FileServer(getFileSystem())

	http.Handle("/foo/count/", http.StripPrefix("/foo/", new(countHandler)))

	// Serve static files
	http.Handle("/app/", http.StripPrefix("/app/", fs))

	log.Println("Listening on http://localhost:3000/")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
