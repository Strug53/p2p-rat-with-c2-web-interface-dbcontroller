package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	//e.Use(middleware.Static())

	/*e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == "123", nil
	}))*/

	e.GET("/ready", func(c echo.Context) error {
		resp, err := http.Get("https://localhost:443/sendClientsTable")
		if err != nil {
			fmt.Printf("Error with post request")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error in parsing of response")
		}
		resp.Body.Close()

		return c.HTML(http.StatusOK, string(body))
	})
	e.GET("/", func(c echo.Context) error {
		return c.File("static/login.html")
	})
	e.POST("/dashboard", func(c echo.Context) error {
		return c.File("static/dashboard.html")
	})
	e.POST("/login", func(c echo.Context) error {
		//usr := c.FormValue("Username")
		//key := c.FormValue("key")
		return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")

	})
	e.Logger.Fatal(e.StartTLS(":8080", "../cert.pem", "../key.pem"))
}
