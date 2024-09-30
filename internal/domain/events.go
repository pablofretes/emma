package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID               uuid.UUID      `gorm:"type:char(36);primaryKey"`
	Title            string         `gorm:"size:255;not null" json:"title" binding:"required,max=255"`
	ShortDescription string         `gorm:"size:500;not null" json:"short_description" binding:"required,max=500"`
	LongDescription  string         `gorm:"type:text" json:"long_description" binding:"max=5000"`
	Date             time.Time      `gorm:"not null" json:"date" binding:"required,future"`
	Location         string         `gorm:"size:255;not null" json:"location" binding:"required,max=255"`
	Status           string         `gorm:"type:enum('draft','published');default:'draft'" json:"status" binding:"oneof=draft published"`
	OrganizerID      uuid.UUID      `gorm:"type:uuid" json:"organizer_id" binding:"required"`
	Attendees        []User         `gorm:"many2many:user_events;" json:"attendees"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type EventRepository interface {
	CreateEvent(event Event) (Event, error)
	GetEvent(id string) (Event, error)
	UpdateEvent(id string, event Event) (Event, error)
	DeleteEvent(id string) error
	GetAllEvents(status, title, date string) ([]Event, error)
}
