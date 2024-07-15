package main

import (
	"dbcontroller"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	//e.Use(middleware.Static())

	/*e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == "123", nil
	}))*/
	clients := dbcontroller.Select_all_clients()

	for _, cc := range clients {
		fmt.Println(cc.IP)
	}
	e.GET("/usr", func(c echo.Context) error {
		return c.String(http.StatusOK, "Yes")
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
	e.Logger.Fatal(e.Start(":8080"))
}
