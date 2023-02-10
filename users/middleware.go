package users

import (
	"go-postgres-fiber/helpers"

	"github.com/gofiber/fiber/v2"
)

func SecureAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		authToken := c.Cookies("AccessToken")
		if authToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"message": "Invalid access token",
			})
		}

		_, claims, err := helpers.VerifyJWT(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"error":   true,
					"message": err.Error(),
				},
			)
		}

		c.Locals("user_id", claims.Issuer)

		return c.Next()

	}
}
