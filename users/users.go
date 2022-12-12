package users

import (
	"github.com/gofiber/fiber/v2"
)

func errorMissingRequiredPayload() error {
	return fiber.NewError(400, "Required fields are missing")
}

func SetupRoutes(app fiber.Router) {
	api := app.Group("/api")
	user := api.Group("/user")

	user.Post("/register", CreateUser)
	user.Post("/login", ValidateUser)
	user.Get("/authenticate", AuthenticateUser)
}
