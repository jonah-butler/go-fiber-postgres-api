// package items

// import (
// 	"fmt"
// 	"net/http"
// 	"time"

// 	database "go-postgres-fiber/connection"
// 	models "go-postgres-fiber/models"

// 	"github.com/gofiber/fiber/v2"
// )

// type Item struct {
// 	Name   string `json:"name"`
// 	UserID uint   `json:"user_id"`
// }

// func SetupRoutes(app fiber.Router) {
// 	api := app.Group("/api")
// 	item := api.Group("/item")

// 	item.Post("/", CreateItem)
// 	item.Get("/", GetItem)
// }

// func CreateItem(context *fiber.Ctx) error {
// 	item := models.ShareableItem{}

// 	err := context.BodyParser(&item)
// 	if err != nil {
// 		return context.Status(http.StatusUnprocessableEntity).JSON(
// 			&fiber.Map{
// 				"status":  400,
// 				"message": "request failed",
// 			},
// 		)
// 	}

// 	item.CreatedAt = time.Now()
// 	item.UpdatedAt = time.Now()

// 	err = database.Conn.Create(&item).Error
// 	if err != nil {
// 		fmt.Println("error creating item", err)
// 		return context.Status(http.StatusBadRequest).JSON(
// 			&fiber.Map{"message": "could not create item"},
// 		)
// 	}

// 	return context.Status(http.StatusOK).JSON(
// 		&fiber.Map{
// 			"status": http.StatusOK,
// 			"item":   item,
// 		},
// 	)

// }

// func GetItem(context *fiber.Ctx) error {

// 	item := models.ShareableItem{}

// 	err := database.Conn.Find(&item).Error
// 	if err != nil {
// 		fmt.Println("could not lookup shareable items")
// 		return context.Status(400).JSON(
// 			&fiber.Map{"message": "could not lookup items"},
// 		)
// 	}

// 	return context.Status(200).JSON(
// 		&fiber.Map{"items": item},
// 	)

// }
