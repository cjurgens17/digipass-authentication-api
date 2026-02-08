package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "" {
		log.Println("Local deploy - manually loading .env file")
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	if err := Run(); err != nil {
		log.Fatal("Application Failed:", err)
	}
}
