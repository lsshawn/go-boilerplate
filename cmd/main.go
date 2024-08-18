package main

import (
	"context"
	"log"

	"boilerplate/internal/database"
	"boilerplate/internal/handlers"
	"boilerplate/views"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app := echo.New()

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

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
