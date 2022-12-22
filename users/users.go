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

	// public user endpoints
	user.Post("/register", CreateUser)
	user.Post("/login", ValidateUser)
	user.Get("/authenticate", AuthenticateUser)
	user.Get("/refresh-access-token", RefreshAccessToken)

	protected := user.Group("/p")
	protected.Use(SecureAuth())

	// private user endpoints
	protected.Get("/private", PrivateUser)
}
