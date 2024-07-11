package main

import (
	"fmt"
	"net/http"
)

func main() {

	fmt.Printf("Starting web.. \n")
	fs := http.FileServer(http.Dir("./static/login.html"))

	http.Handle("/", fs)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.FileServer(http.Dir("./static/login.html"))
	// })
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Query().Get("m"))
	})

	http.ListenAndServe(":8080", nil)
}
