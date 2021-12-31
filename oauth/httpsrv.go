package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "---BEGIN HEADERS---\n")
	for k, v := range r.Header {
		fmt.Fprintf(w, " %-30s = %q\n", k, v)
	}
	fmt.Fprint(w, "---END HEADERS---\n")

	var d struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello World!")
		return
	}
	if d.Message == "" {
		fmt.Fprint(w, "Hello World!")
		return
	}
	fmt.Fprint(w, html.EscapeString(d.Message))
}

func main() {
	http.HandleFunc("/", HelloWorld)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
