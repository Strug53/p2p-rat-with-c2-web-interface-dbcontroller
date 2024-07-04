package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type readyStruct struct {
	ip string
}

func main() {
	fmt.Printf("Starting..")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Query().Get("m"))
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Printf("Check ready ")
		rdy := &readyStruct{}
		err := json.NewDecoder(r.Body).Decode(rdy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf(rdy.ip)
	})

	http.ListenAndServe(":1234", nil)
}
