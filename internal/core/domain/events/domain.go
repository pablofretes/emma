package events

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title            string    `gorm:"size:255;not null" json:"title" binding:"required,max=255"`
	ShortDescription string    `gorm:"size:500;not null" json:"short_description" binding:"required,max=500"`
	LongDescription  string    `gorm:"type:text" json:"long_description" binding:"max=5000"`
	Date             time.Time `gorm:"not null" json:"date" binding:"required,future"`
	Location         string    `gorm:"size:255;not null" json:"location" binding:"required,max=255"`
	Status           string    `gorm:"type:enum('draft','published');default:'draft'" json:"status" binding:"oneof=draft published"`
	OrganizerID      uint      `gorm:"foreignKey:ID" json:"organizer_id" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type EventService interface {
	CreateEvent(event Event) (Event, error)
	GetEvent(id string) (Event, error)
	UpdateEvent(id string, event Event) (Event, error)
	DeleteEvent(id string) error
	GetAllEvents() ([]Event, error)
}
