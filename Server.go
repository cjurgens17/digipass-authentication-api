package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/joho/godotenv"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. If running Local add .env or else railway variables are being used")
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