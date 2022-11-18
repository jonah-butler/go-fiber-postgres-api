package users

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	database "go-postgres-fiber/connection"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type User struct {
	Id        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UnvalidatedUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func errorMissingRequiredPayload() error {
	return fiber.NewError(400, fmt.Sprintf("Required fields are missing"))
}

func SetupRoutes(app fiber.Router) {
	api := app.Group("/api")
	user := api.Group("/user")

	user.Post("/register", CreateUser)
	user.Post("/login", ValidateUser)
}

func CreateUser(context *fiber.Ctx) error {
	user := User{}

	err := context.BodyParser(&user)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"},
		)
		return err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Password = hasAndSaltPassword(user.Password)

	err = database.Conn.Create(&user).Error
	if err != nil {
		fmt.Println("error creating user", err)
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"},
		)
	}

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status": http.StatusOK,
			"user":   user,
		},
	)
}

func ValidateUser(context *fiber.Ctx) error {

	userPayload := UnvalidatedUser{}
	user := User{}

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

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status": http.StatusOK,
			"user":   user,
		},
	)

}

func hasAndSaltPassword(password string) string {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}

func validatePassword(dbPassword string, plainPwd string) bool {
	byteHash := []byte(dbPassword)
	if err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd)); err != nil {
		return false
	}
	return true
}
