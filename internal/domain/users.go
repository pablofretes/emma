package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              uuid.UUID      `gorm:"type:char(36);primaryKey"`
	Username        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"`
	Password        string         `gorm:"type:varchar(255);not null" json:"-"`
	Role            string         `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	OrganizedEvents []Event        `gorm:"foreignKey:OrganizerID" json:"organized_events"`
	AttendingEvents []Event        `gorm:"many2many:user_events;" json:"attending_events"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRepository interface {
	Create(user User) error
	GetById(id string) (User, error)
	GetByUsername(username string) (User, error)
	Update(eventID string, userID string) error
	Login(username string, password string) (User, error)
	GetUsersEvents(userID string, status string) ([]Event, error)
}
