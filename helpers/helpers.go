package helpers

import (
	"fmt"
	database "go-postgres-fiber/connection"
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
func GenerateAccessClaims(user models.User) (*models.JWTClaims, string) {

	claim := &models.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    user.ID.String(),
			ExpiresAt: generateJWTExp(15),
			Subject:   "access_token",
			IssuedAt:  time.Now().Unix(),
		},
		Email:    user.Email,
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(claim)
	return claim, tokenString

}

func GenerateRefreshClaims(claims *models.JWTClaims) string {

	// check if claims issuer has any refresh tokens stored in db
	result := database.Conn.Where(&models.JWTRefreshClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: claims.Issuer,
		},
	}).Find(&models.JWTRefreshClaims{})

	// if refresh token already present, delete before inserting new claim
	if result.RowsAffected == 1 {
		database.Conn.Where(&models.JWTRefreshClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer: claims.Issuer,
			},
		}).Delete(&models.JWTRefreshClaims{})
	}

	refreshClaim := &models.JWTRefreshClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    claims.Issuer,
			ExpiresAt: generateJWTRefreshExp(15),
			Subject:   "refresh_token",
			IssuedAt:  time.Now().Unix(),
		},
	}
	// create the new claim in db
	database.Conn.Create(&refreshClaim)

	// create new jwt
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
	refreshTokenStr, err := refreshToken.SignedString(secretKey)
	if err != nil {
		panic(err)
	}

	return refreshTokenStr

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
func VerifyJWT(headers map[string]string) (*jwt.Token, *models.JWTClaims, error) {
	auth := headers["Authorization"]

	preToken := strings.Split(auth, " ")[1]
	claims := new(models.JWTClaims)

	if len(preToken) > 0 {

		token, err := jwt.ParseWithClaims(preToken, claims,
			func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

		if err != nil {
			panic("error occurred parsing auth token")
		}
		if token.Valid {
			return token, claims, nil
		}

	}

	return nil, nil, fmt.Errorf("unauthorized access")
}

func generateJWTRefreshExp(days int) int64 {
	return time.Now().Add((time.Hour * 24) * time.Duration(days)).Unix()
}

func generateJWTExp(minutes int) int64 {
	minutesConverted := time.Duration(minutes) * time.Minute
	return time.Now().Add(time.Minute * minutesConverted).Unix()
}
