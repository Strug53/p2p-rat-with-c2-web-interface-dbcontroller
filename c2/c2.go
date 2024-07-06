package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type readyStruct struct {
	IP string `json:"ip"`
}

func setCommand(rdy *readyStruct) {
	fmt.Printf("Enter the command here: ")

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	url := "http://" + rdy.IP + "/cmd?m=" + text
	fmt.Println(text)
	fmt.Printf("\n")

	resp, _ := http.Post(url, "application/json", nil) //edit exceptions
	resp.Body.Close()

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
		setCommand(&rdy)
	})

	http.ListenAndServe(":1234", nil)
}
