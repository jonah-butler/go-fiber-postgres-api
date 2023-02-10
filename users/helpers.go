package users

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func hashAndSaltPassword(password string) string {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
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

func GetAuthCookies(accessToken, refreshToken string) (*fiber.Cookie, *fiber.Cookie) {

	accessTokenCookie := &fiber.Cookie{
		Name:     "AccessToken",
		Value:    accessToken,
		HTTPOnly: false,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
	}

	refreshTokenCookie := &fiber.Cookie{
		Name:     "RefreshToken",
		Value:    accessToken,
		HTTPOnly: false,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		Secure:   false,
	}

	return accessTokenCookie, refreshTokenCookie

}
