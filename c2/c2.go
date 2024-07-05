package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type readyStruct struct {
	IP string `json:"ip"`
}

func main() {
	fmt.Printf("Starting.. \n")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Query().Get("m"))
	})
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("ready endpoint used \n")
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var rdy readyStruct
		err := json.NewDecoder(r.Body).Decode(&rdy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Ip: " + rdy.IP)
	})

	http.ListenAndServe(":1234", nil)
}
