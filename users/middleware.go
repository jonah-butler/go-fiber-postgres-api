package users

import (
	"go-postgres-fiber/helpers"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func SecureAuth() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		_, claims, err := helpers.VerifyJWT(c.GetReqHeaders())
		if err != nil {
			fiber.NewError(http.StatusUnauthorized, err.Error())
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
			return fiber.NewError(http.StatusUnauthorized, "JWT validation error")
		}

		c.Locals("id", claims.Issuer)

		return c.Next()

	}
}
