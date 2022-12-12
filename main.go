package main

import (
	"log"
	"os"

	database "go-postgres-fiber/connection"
	// "go-postgres-fiber/items"
	"go-postgres-fiber/models"
	"go-postgres-fiber/users"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.NewConnection(config)
	if err != nil {
		log.Fatal("could not connect to databse")
	}

	err = models.MigrateTables(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	database.Conn = db

	app := fiber.New()
	users.SetupRoutes(app)
	// items.SetupRoutes(app)
	app.Listen(":8080")
}
