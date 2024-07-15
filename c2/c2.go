package main

import (
	"bufio"
	"dbcontroller"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
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
type Answer struct {
	IP          string `json:"Ip"`
	Command     string `json:"Command"`
	Key_Manager string `json:"Key_Manager"`
	Key_Client  string `json:"Key_Client"`
	Result      string `json:"Result"`
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
	e := echo.New()
	fmt.Printf("Starting.. \n")
	e.GET("/sendClientsTable", func(c echo.Context) error {
		var html_string string

		clients := dbcontroller.Select_all_clients()
		fmt.Println(clients)
		fmt.Printf("\n")
		for _, c := range clients {
			html_buf := fmt.Sprintf(`
			<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
		  	</tr>`, strconv.Itoa(c.Id), c.IP, c.Port, c.System, c.Client_key, c.Date)

			html_string += html_buf
		}

		return c.HTML(http.StatusOK, html_string)
	})
	e.POST("/answer", func(c echo.Context) error {
		fmt.Printf("\n")

		var Answer Answer
		err := json.NewDecoder(c.Request().Body).Decode(&Answer)
		if err != nil {
			fmt.Printf(err.Error())
		}
		fmt.Printf("Answer coming from client: %s \n", Answer.IP)
		fmt.Printf("\t %s \n", Answer.Command)
		fmt.Printf("\t %s \n", Answer.Key_Manager)
		fmt.Printf("\t %s \n", Answer.Key_Client)
		fmt.Printf("\t %s \n", Answer.Result)

		fmt.Printf("\n")
		return c.String(http.StatusOK, "Ok")

	})
	e.POST("/ready", func(c echo.Context) error {
		var Contact Contact
		err := json.NewDecoder(c.Request().Body).Decode(&Contact)
		if err != nil {
			fmt.Printf(err.Error())
		}
		fmt.Printf("Client ready for commands:	\n")
		fmt.Printf("\t %s \n", Contact.IP)
		fmt.Printf("\t %s \n", Contact.Port)
		fmt.Printf("\t %s \n", Contact.System)
		fmt.Printf("\t %s \n", Contact.Date)

		fmt.Printf("\n")
		go setCommand()
		return c.String(http.StatusOK, "Ok")
		//ip = rdy.IP
	})

	e.Logger.Fatal(e.Start(":1234"))
}
