package main

import (
	"log"

	"boilerplate/internal/database"

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

	app.Static("/static", "static")

	app.Logger.Fatal(app.Start(":1323"))
}
