package main

import (
	"DigiPassAuthenticationApi/routes"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"os"
)

//ValidatorWrapper

type ValidatorWrapper struct {
	validator *validator.Validate
}

func (vw *ValidatorWrapper) Validate(i interface{}) error {
	return vw.validator.Struct(i)
}

func Run() error {
	db := initDB()
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// Middleware to Inject DB into each request context
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	//Set Up Validator Wrapper
	e.Validator = &ValidatorWrapper{validator: validator.New()}


	//need to pass db connection to handlers, or service layer
	routes.SetUpRoutes(e)
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
