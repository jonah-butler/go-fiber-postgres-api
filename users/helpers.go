package users

import (
	"fmt"

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
