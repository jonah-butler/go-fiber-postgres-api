package models

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"primary_key; unique; not null; type:uuid; column:id; default:gen_random_uuid()" json:"id"`
	Username  string         `gorm:"unique; not null" json:"username"`
	Email     string         `gorm:"unique" json:"email"`
	Password  string         `gorm:"type:varchar(500)" json:"password"`
	UpdatedAt time.Time      `gorm:"type:time" json:"updated_at"`
	CreatedAt time.Time      `gorm:"type:time" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type MinimumUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type JWTClaims struct {
	jwt.StandardClaims
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type JWTRefreshClaims struct {
	jwt.StandardClaims
	ID string `gorm:"primaryKey"`
}

type UserErrors struct {
	Err      bool   `json:"error"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func MigrateTables(db *gorm.DB) error {

	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("failed to migrate user model")
		return err
	}
	err = db.AutoMigrate(&JWTRefreshClaims{})
	if err != nil {
		log.Fatal("failed to migrate location model")
		return err
	}
	return err
}
