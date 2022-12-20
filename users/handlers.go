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

	token, _, err := helpers.VerifyJWT(context.GetReqHeaders())
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

	claim, token := helpers.GenerateAccessClaims(user)
	refreshToken := helpers.GenerateRefreshClaims(claim)
	fmt.Println("subject = ", claim.Subject)
	fmt.Println("issuer = ", claim.Issuer)
	// fmt.Println("subject = ", claim.Subject)

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
