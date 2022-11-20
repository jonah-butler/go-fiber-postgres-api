package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	Longitude *string `gorm:"type:varchar(50)" json:"longitude"`
	Latitude  *string `gorm:"type:varchar(50)" json:"latitude"`
}

type ShareableItem struct {
	gorm.Model
	ID     uint    `gorm:"primary key,autoIncrement" json:"id"`
	Name   *string `gorm:"type:text" json:"name"`
	User   *User   `gorm:"-"`
	UserID uint    `json:"user_id"`
}

type User struct {
	gorm.Model
	ID             uint            `gorm:"primary key,autoIncrement" json:"id"`
	Username       *string         `gorm:"unique" json:"username"`
	Email          *string         `gorm:"unique" json:"email"`
	Password       *string         `gorm:"type:varchar(500)" json:"password"`
	Location       Location        `gorm:"-" json:"location"`
	UpdatedAt      time.Time       `gorm:"type:time" json:"updated_at"`
	CreatedAt      time.Time       `gorm:"type:time" json:"created_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
	ShareableItems []ShareableItem `gorm:"ForeignKey:UserID"`
}

func MigrateTables(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("failed to migrate user model")
		return err
	}
	err = db.AutoMigrate(&Location{})
	if err != nil {
		log.Fatal("failed to migrate location model")
		return err
	}
	err = db.AutoMigrate(&ShareableItem{})
	if err != nil {
		log.Fatal("failed to migrate shareable items model")
		return err
	}
	return err
}
