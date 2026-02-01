package main

import (
	"DigiPassAuthenticationApi/routes"
	"log"
	"os"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() error {
	db := initDB()
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	//need to pass db connection to handlers, or service layer
	routes.SetUpRoutes(e, db)
	return e.Start(":1323")
}

func initDB() *gorm.DB {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL))
	if err != nil {
		log.Fatalf("Unable to connect database : %v", err)
	}

	log.Println("Connected to Database")
	return db
}
