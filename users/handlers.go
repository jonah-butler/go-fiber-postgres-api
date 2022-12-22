package users

import (
	"fmt"
	database "go-postgres-fiber/connection"
	"go-postgres-fiber/helpers"
	"go-postgres-fiber/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateUser(context *fiber.Ctx) error {

	headers := context.GetReqHeaders()
	authToken := headers["Authorization"]

	token, _, err := helpers.VerifyJWT(authToken)
	if err != nil {
		fiber.NewError(http.StatusUnauthorized, err.Error())
	}

	return context.Status(http.StatusBadRequest).JSON(
		&fiber.Map{
			"message": "authenticated",
			"user":    token,
		},
	)

}

func CreateUser(context *fiber.Ctx) error {
	user := models.User{}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(422).JSON(
			&fiber.Map{"error": "request failed"},
		)
		return err
	}

	if user.Password == "" || user.Email == "" || user.Username == "" {
		return context.Status(400).JSON(
			&fiber.Map{"error": "missing required field(s)"},
		)
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Password = hashAndSaltPassword(user.Password)

	err = database.Conn.Create(&user).Error
	if err != nil {
		fmt.Println("error creating user: ", err)
		return context.Status(500).JSON(
			&fiber.Map{"error": err.Error()},
		)
	}

	newUser := NewUser{}
	newUser.Email = user.Email
	newUser.Username = user.Username

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{"user": newUser},
	)
}

func ValidateUser(context *fiber.Ctx) error {

	userPayload := UnvalidatedUser{}
	user := models.User{}

	if err := context.BodyParser(&userPayload); err != nil {
		return err
	}

	if userPayload.Email == "" || userPayload.Password == "" {
		return errorMissingRequiredPayload()
	}

	database.Conn.
		Where("email = ?", userPayload.Email).First(&user)

	fmt.Println(user)

	isPasswordCorrect := validatePassword(user.Password, userPayload.Password)

	if !isPasswordCorrect {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"status":  http.StatusBadRequest,
				"message": "The provided login details are incorrect",
			},
		)
	}

	claim, token, err := helpers.GenerateAccessClaims(user)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"error":   true,
				"message": err.Error(),
			},
		)
	}
	refreshToken, err := helpers.GenerateRefreshClaims(claim)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"error":   true,
				"message": err.Error(),
			},
		)
	}

	cookie1 := new(fiber.Cookie)
	cookie2 := new(fiber.Cookie)

	cookie1.Name = "AccessToken"
	cookie1.Value = token

	cookie2.Name = "RefreshToken"
	cookie2.Value = refreshToken

	context.Cookie(cookie1)
	context.Cookie(cookie2)

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status":       http.StatusOK,
			"token":        token,
			"refreshToken": refreshToken,
		},
	)

}

func PrivateUser(context *fiber.Ctx) error {

	fmt.Println("made it to private user route")

	return nil

}

func RefreshAccessToken(context *fiber.Ctx) error {

	// check cookie for refresh token
	refreshToken := context.Cookies("RefreshToken")
	if refreshToken == "" {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"error":   true,
				"message": "Request missing required cookie",
			},
		)
	}

	// parse refresh token
	token, claims, err := helpers.VerifyRefreshJWT(refreshToken)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"errror":  true,
				"message": err.Error(),
			},
		)
	}

	// lookup refresh token from cookie
	result := database.Conn.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ?",
		claims.ExpiresAt, claims.IssuedAt, claims.Issuer,
	).First(&models.JWTRefreshClaims{})

	// if no data from table, clear cookies as they've been malformed
	if result.RowsAffected <= 0 {
		context.ClearCookie("AccessToken", "RefreshToken")
		return context.Status(http.StatusForbidden).JSON(
			&fiber.Map{
				"error":   true,
				"message": "Invalid token",
			},
		)
	}

	if token.Valid {

		if claims.ExpiresAt < time.Now().Unix() {
			context.ClearCookie("AccessToken", "RefreshToken")
			return context.Status(http.StatusForbidden).JSON(
				&fiber.Map{
					"error":   true,
					"message": "Invalid token",
				},
			)
		}

	} else {

		context.ClearCookie("AccessToken", "RefreshToken")
		return context.Status(http.StatusForbidden).JSON(
			&fiber.Map{
				"error":   true,
				"message": "Invalid token",
			},
		)

	}

	user := models.User{}

	database.Conn.
		Where("id = ?", claims.Issuer).First(&user)

	_, newToken, err := helpers.GenerateAccessClaims(user)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"error":   true,
				"message": err.Error(),
			},
		)
	}

	context.Cookie(
		&fiber.Cookie{
			Name:     "AccessToken",
			Value:    newToken,
			HTTPOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
			Secure:   true,
		},
	)

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status": http.StatusOK,
			"token":  newToken,
		},
	)

}
