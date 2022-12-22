package users

import (
	"go-postgres-fiber/helpers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func SecureAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		headers := c.GetReqHeaders()
		authToken := headers["Authorization"]

		_, claims, err := helpers.VerifyJWT(authToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"error":   true,
					"message": "Token Expired",
				},
			)
		}
		if claims.ExpiresAt < time.Now().Unix() {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"error":   true,
					"message": "Token Expired",
				},
			)
		}

		ve, _ := err.(*jwt.ValidationError)
		if ve != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"error":   true,
					"message": "Token Validation Error",
				},
			)
		}

		c.Locals("id", claims.Issuer)

		return c.Next()

	}
}
