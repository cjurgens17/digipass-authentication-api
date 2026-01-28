package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)


func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "" {
		log.Println("Local deploy - manually loading .env file")
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	ctx := context.Background()
	pool := NewDBPool(ctx)

	app := &App{
		DB: pool,
	}
	defer app.DB.Close()

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}