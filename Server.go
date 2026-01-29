package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
    DB *pgxpool.Pool
}

func NewDBPool(ctx context.Context) *pgxpool.Pool{
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to conect database: %v", err)
	}

	var now string
	if err := pool.QueryRow(ctx, "SELECT NOW()::text").Scan(&now); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Println("Connected to database at", now)
	return pool
}


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