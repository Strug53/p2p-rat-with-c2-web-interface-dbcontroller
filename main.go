package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Query().Get("m"))
	})

	http.ListenAndServe(":1234", nil)
}
