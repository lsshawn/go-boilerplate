package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/lsshawn/go-todo/views"
)

func main() {
	app := echo.New()

	app.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	app.GET("/", func(c echo.Context) error {
		component := views.Index()
		return component.Render(context.Background(), c.Response().Writer)
	})

	app.Static("/css", "css")
	app.Static("/static", "static")

	app.Logger.Fatal(app.Start(":1323"))
}
