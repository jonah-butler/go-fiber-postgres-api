package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	database "go-postgres-fiber/connection"
	"go-postgres-fiber/helpers"
	"go-postgres-fiber/models"
	"go-postgres-fiber/users"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type User struct {
	Username string ``
	Email    string ``
	Password string ``
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"},
		)
		return err
	}

	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"},
		)
		return err
	}

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "new book added"},
	)
}

func (r *Repository) CreateJWT(context *fiber.Ctx) error {

	token, err := helpers.GenerateJWT()
	if err != nil {
		log.Fatal("failed to execute GenerateJWT func", err)
	}

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{"token": token},
	)

}

func (r *Repository) ValidateJWT(context *fiber.Ctx) error {

	helpers.VerifyJWT(context.GetReqHeaders())

	return context.Status(http.StatusAccepted).JSON(
		&fiber.Map{"status": 200},
	)

}

// func (r *Repository) GetBooks(context *fiber.Ctx) error {
// 	books := &[]models.Books{}

// 	err := r.DB.Find(books).Error
// 	if err != nil {
// 		context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not retrieve books"},
// 		)
// 		return err
// 	}

// 	return context.Status(http.StatusOK).JSON(
// 		&fiber.Map{"data": books},
// 	)
// }

func DeleteBook() {
	fmt.Println("delete book")
}

func GetBookByID() {
	fmt.Println("search by i")
}

func UpdateBook() {
	fmt.Println(" update book")
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	// api.Get("/books", r.GetBooks)
	api.Get("/jwt", r.CreateJWT)
	api.Post("/jwt", r.ValidateJWT)
	api.Post("/register")
}

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

	// r := Repository{
	// 	DB: db,
	// }

	app := fiber.New()
	users.SetupRoutes(app)
	app.Listen(":8080")
}
