package users

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	database "go-postgres-fiber/connection"
	"go-postgres-fiber/helpers"
	"go-postgres-fiber/models"

	"github.com/gofiber/fiber/v2"
)

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
	user.Get("/", GetUser)
	user.Get("/authenticate", AuthenticateUser)
}

type TokenPayload struct {
	Token string `json:"token"`
}

func AuthenticateUser(context *fiber.Ctx) error {

	// token := TokenPayload{}

	token := helpers.VerifyJWT(context.GetReqHeaders())
	if token != nil {
		fmt.Println("token validated", token)
	}
	// if err != nil {
	// 	fmt.Println("error parsing token payload")
	// }

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
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"},
		)
		return err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	fmt.Println("before hasing", user.Password)
	user.Password = hashAndSaltPassword(user.Password)
	fmt.Println("after hasing", user.Password)

	err = database.Conn.Create(&user).Error
	if err != nil {
		fmt.Println("error creating user", err)
		return context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create user"},
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

	token, err := helpers.GenerateJWT(user)
	if err != nil {
		fmt.Println("failed to create token")
	}

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status": http.StatusOK,
			"token":  token,
		},
	)

}

func GetUser(context *fiber.Ctx) error {

	user := models.User{}

	err := database.Conn.
		Preload("ShareableItems").
		First(&user, 5).Error

	if err != nil {
		fmt.Println("failed to lookup user of ID 5")
	}

	return context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"status": http.StatusOK,
			"user":   user,
		},
	)
}

func hashAndSaltPassword(password string) string {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}

func validatePassword(dbPassword string, plainPwd string) bool {
	byteHash := []byte(dbPassword)
	if err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd)); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
