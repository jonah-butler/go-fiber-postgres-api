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

	token := helpers.VerifyJWT(context.GetReqHeaders())
	if token != nil {
		fmt.Println("token validated", token)
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

	isPasswordCorrect := validatePassword(user.Password, userPayload.Password)

	if !isPasswordCorrect {
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"status":  http.StatusBadRequest,
				"message": "The provided login details are incorrect",
			},
		)
	}

	// token, err := helpers.GenerateJWT(user)
	// if err != nil {
	// 	fmt.Println("failed to create token")
	// }
	claim, token := helpers.GenerateAccessClaims(user.ID.String())
	fmt.Println(claim, token)

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status": http.StatusOK,
			"token":  token,
		},
	)

}
