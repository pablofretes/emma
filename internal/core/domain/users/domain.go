package users

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Role    string `gorm:"type:enum('admin','user');default:'user'" json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserService interface {
	CreateUser(user User) error
	GetUserByUsername(username string) (User, error)
}

