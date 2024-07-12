package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	//e.Use(middleware.Static())

	/*e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == "123", nil
	}))*/

	e.GET("/", func(c echo.Context) error {
		return c.File("static/login.html")
	})
	e.GET("/dashboard", func(c echo.Context) error {
		return c.File("static/dashboard.html")
	})
	e.POST("/login", func(c echo.Context) error {
		//usr := c.FormValue("Username")
		//key := c.FormValue("key")
		return c.Redirect(http.StatusTemporaryRedirect, "/dashboard")

	})
	e.Logger.Fatal(e.Start(":8080"))
}
