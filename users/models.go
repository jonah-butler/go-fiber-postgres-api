package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserAuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
