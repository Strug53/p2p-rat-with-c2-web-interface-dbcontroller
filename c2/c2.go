package main

import (
	"dbcontroller"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var ip = "localhost:5555"

var timeOfRefresh int = 50

type Contact struct {
	IP     string `json:"Ip"`
	Port   string `json:"Port"`
	System string `json:"System"` //os
	Key    string `json:"Key"`    //new programm -> new individual key for autorizaiton (make with time and load main function)
	//manager key also
	Date string `json:"Date"`
}

// for logging
type Answer struct {
	IP          string `json:"Ip"`
	Command     string `json:"Command"`
	Key_Manager string `json:"Key_Manager"`
	Key_Client  string `json:"Key_Client"`
	Result      string `json:"Result"`
}
type cmdForm struct {
	ID  string `json:"id"`
	Cmd string `json:"cmd"`
}

var secretKey string = "Bl5cMvdxpej8efG3vCQBl0UVLFByoQ9W"

func sendCommand(c echo.Context) error {

	fmt.Printf("Enter the command here: ")
	cmdJson := cmdForm{}

	err := json.NewDecoder(c.Request().Body).Decode(&cmdJson)

	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("\n")

	cmdJson.Cmd = strings.Replace(cmdJson.Cmd, " ", "+", -1)
	fmt.Printf(cmdJson.Cmd)
	url := "http://" + ip + "/cmd?m=" + cmdJson.Cmd
	fmt.Printf(url)

	fmt.Printf("\n")

	fmt.Printf(cmdJson.ID)

	resp, err := http.Post(url, "application/json", nil) //edit exceptions
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ans := Answer{}
	er := json.NewDecoder(resp.Body).Decode(&ans)
	if er != nil {
		fmt.Printf("Error in decoding answer")
	}
	fmt.Println(ans)
	return c.String(http.StatusOK, ans.Result)
}
func refreshDatabase() {

	time1 := time.NewTimer(time.Duration(timeOfRefresh) * time.Second)
	<-time1.C

	fmt.Printf("\nREFRESH\n")

	dbcontroller.Delete_all()

	refreshDatabase()
}

func main() {
	go refreshDatabase()
	startServer()
}
func startServer() {
	e := echo.New()

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:auth-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			fmt.Println(key)
			return key == secretKey, nil
		},
	}))

	fmt.Printf("Starting.. \n")
	//htmx
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
				<td>
				<a href="https://localhost:7777/remote?uid=%s">К управлению</a>
				</td>
		  	</tr>`, strconv.Itoa(c.Id), c.IP, c.Port, c.System, c.Client_key, c.Date, strconv.Itoa(c.Id))

			html_string += html_buf
		}

		return c.HTML(http.StatusOK, html_string)
	})
	//getCommand
	e.POST("/getCommand", sendCommand)

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
		cont := dbcontroller.Select_client_IP(Contact.IP)
		fmt.Println(Contact)
		fmt.Printf("\n")

		fmt.Println(cont)
		//Change ip to key
		if cont.IP == Contact.IP {
			return c.String(http.StatusOK, "Ok")
		}
		dbcontroller.Add_client(Contact.IP, Contact.Port, Contact.System, Contact.Key, Contact.Date)
		return c.String(http.StatusOK, "Ok")
		//ip = rdy.IP
	})
	//adding users in db
	//SECURITY!!. ADD VERIFICATION
	e.GET("/checkKey", func(c echo.Context) error {
		fmt.Printf("\nCheckKEy\n")
		return c.String(http.StatusOK, "Ok")
	})

	e.GET("/checkUser/:uid", func(c echo.Context) error {
		userId := c.Param("uid")
		fmt.Printf("\n")
		fmt.Printf(userId)
		fmt.Printf("\n")

		userId_int, err := strconv.Atoi(userId)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "None")
		}
		db_result := dbcontroller.Select_client_ID(userId_int)
		fmt.Println(db_result)
		if (db_result != dbcontroller.Client{}) {
			return c.String(http.StatusOK, "Ok")
		} else {
			return c.String(http.StatusOK, "None")
		}
	})
	e.Logger.Fatal(e.StartTLS(":443", "../cert.pem", "../key.pem"))
}
