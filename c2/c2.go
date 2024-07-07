package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var ip = "localhost:5555"

type readyStruct struct {
	IP string `json:"ip"`
}

func setCommand() {

	for {
		if ip == "" {
			continue
		}

		fmt.Printf("Enter the command here: ")
		var command string
		fmt.Fscan(os.Stdin, &command) // fix the problem with space

		url := "http://" + ip + "/cmd?m=" + command
		fmt.Printf(url)

		fmt.Printf("\n")

		resp, err := http.Post(url, "application/json", nil) //edit exceptions
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

func main() {
	startServer()
}
func startServer() {
	fmt.Printf("Starting.. \n")
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

		var rdy readyStruct
		err := json.NewDecoder(r.Body).Decode(&rdy)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Client ready for commands:	" + rdy.IP)
		fmt.Printf("\n")
		go setCommand()
		//ip = rdy.IP
	})

	http.ListenAndServe(":1234", nil)
}
