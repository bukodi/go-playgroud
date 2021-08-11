package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
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

func main() {
	fs := http.FileServer(getFileSystem())

	// Serve static files
	http.Handle("/", fs)

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
