package users

import (
	"time"
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

type NewUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
