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

/*
*
* Handler for creating a new user
*
 */
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

/*
*
* Handler for processing user payload when logging in
* - process payload
* - perform lookup and password validation
* - generate jwt access claims
* - generate jwt refresh claims
* - generate and set access/refresh cookies
*
*/
func ValidateUser(context *fiber.Ctx) error {

	userPayload := UserAuthPayload{}
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
				"message": fmt.Sprintf("The provided login details are incorrect for email: %s", userPayload.Email),
			},
		)
	}

	claim, access_token, err := helpers.GenerateAccessClaims(user)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"error":   true,
				"message": err.Error(),
			},
		)
	}
	refresh_token, err := helpers.GenerateRefreshClaims(claim)
	if err != nil {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"error":   true,
				"message": err.Error(),
			},
		)
	}

	accessCookie, refreshCookie := GetAuthCookies(access_token, refresh_token)

	context.Cookie(accessCookie)
	context.Cookie(refreshCookie)

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status": http.StatusOK,
			"token":        access_token,
			"refresh_token": refresh_token,
		},
	)

}

/*
*
* Handler for authenticating requests
*
*/
func AuthenticateUser(context *fiber.Ctx) error {

	authToken := context.Cookies("AccessToken")
	if authToken == "" {
		return fiber.NewError(http.StatusBadRequest, "Invalid token")
	}

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

func RefreshAccessToken(context *fiber.Ctx) error {

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

func GetUser(context *fiber.Ctx) error {

	userId := context.Locals("user_id")
	// check user id from jwt parse is an actual value
	if userId == "" || userId == nil {
		return context.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error":   true,
				"message": "Invalid user id from jwt parse",
			},
		)
	}

	// if user id from locals does not match user id param, access to user data is not allowed
	if userId != context.Params("id") {
		return context.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error":   true,
				"message": "Access unauthorized",
			},
		)
	}

	user := User{}

	database.Conn.
		Where("id = ?", userId).First(&user)

	return context.Status(200).JSON(
		fiber.Map{
			"user": user,
		},
	)
}
