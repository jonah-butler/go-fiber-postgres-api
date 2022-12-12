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

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// that new new - for creating a more robust jwt claims
func GenerateAccessClaims(uuid string) (*models.JWTClaims, string) {

	claim := &models.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    uuid,
			ExpiresAt: generateJWTExp(15),
			Subject:   "access_token",
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return claim, tokenString

}

// this helper generates a JWT used for validation across services
// this jwt generation is strictly for short lived authorization tokens
func GenerateJWT(user models.User) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("failed to load env vars for jwt token generation", err)
	}

	// generate token with signing method
	token := jwt.New(jwt.SigningMethodHS256)

	// modify jwt via Claims method
	claims := token.Claims.(jwt.MapClaims)
	claims["expiration"] = generateJWTExp(15)
	claims["user"] = user

	tokenString, err := token.SignedString(secretKey)
	fmt.Println("token string", tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func GenerateRefreshJWT(user models.User) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("failed to load env vars for jwt token generation", err)
	}

	var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	// generate token with signing method
	token := jwt.New(jwt.SigningMethodHS256)

	// modify jwt via Claims method
	claims := token.Claims.(jwt.MapClaims)
	claims["expiration"] = generateJWTRefreshExp(60)
	claims["id"] = user.ID

	tokenString, err := token.SignedString(secretKey)

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

func generateJWTRefreshExp(days int) time.Time {
	return time.Now().Add((time.Hour * 24) * time.Duration(days))
}

func generateJWTExp(minutes int) int64 {
	minutesConverted := time.Duration(minutes) * time.Minute
	return time.Now().Add(time.Minute * minutesConverted).Unix()
}
