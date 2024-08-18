package main

import (
	"context"

	"boilerplate/internal/handlers"
	"boilerplate/views"

	"github.com/labstack/echo/v4"
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

	app.GET("/account", handlers.AccountPage)
	app.POST("/account/request-otp", handlers.RequestOTP)
	app.POST("/account/validate-otp", handlers.ValidateOTP)
	app.POST("/account/logout", handlers.Logout)

	app.Logger.Fatal(app.Start(":1323"))
}
