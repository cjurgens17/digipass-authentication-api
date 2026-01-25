package main

import (
	"os"
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
    DB *pgxpool.Pool
}

func NewDBPool(ctx context.Context) *pgxpool.Pool{
	dbURL := os.Getenv("PUBLIC_DATABASE_URL")
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