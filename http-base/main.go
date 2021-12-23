package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	server := gee.New()
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL)
	})
	server.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header [%q] = %q\n", k, v)
		}
	})
	server.Run(":8080")
}
