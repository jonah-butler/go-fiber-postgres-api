package helpers

import (
	"fmt"
	"go-postgres-fiber/models"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func GenerateJWT(user models.User) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("failed to load env vars for jwt token generation", err)
	}

	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	// generate token with signing method
	token := jwt.New(jwt.SigningMethodHS256)

	// modify jwt via Claims method
	claims := token.Claims.(jwt.MapClaims)
	claims["expiration"] = generateJWTExp(7)
	claims["user"] = user

	tokenString, err := token.SignedString(secretKey)
	fmt.Println("token string", tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// func VerifyJWT(headers map[string]string) {
func VerifyJWT(headers map[string]string) *jwt.Token {
	auth := headers["Authorization"]

	preToken := strings.Split(auth, " ")[1]

	if len(preToken) > 0 {

		token, err := jwt.Parse(preToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			log.Fatal("error occurred parsing auth token")
		}
		if token.Valid {
			// fmt.Println("valid token", token)
			return token
		}

	}
	return nil
}

func generateJWTExp(days int) time.Time {
	return time.Now().Add((time.Hour * 24) * time.Duration(days))
}
