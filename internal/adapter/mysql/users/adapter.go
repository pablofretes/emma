package adapters

import (
	"emma/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type MySQLAdapter struct {
	DB *gorm.DB
}

func NewMySQLAdapter(db *gorm.DB) (*MySQLAdapter) {
	return &MySQLAdapter{DB: db}
}

func (s *MySQLAdapter) Create(user domain.User) error {
	user.ID = uuid.New()
	result := s.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *MySQLAdapter) GetById(id string) (domain.User, error) {
	var user domain.User
	result := s.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (s *MySQLAdapter) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	result := s.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (s *MySQLAdapter) Login(username string, password string) (domain.User, error) {
	var user domain.User
	result := s.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	return user, nil
}

func (s *MySQLAdapter) Update(eventID string, userID string) error {
	var event domain.Event
	if err := s.DB.First(&event, "id = ?", eventID).Error; err != nil {
		return err
	}

	if event.Date.Before(time.Now()) {
		return fmt.Errorf("event date is in the past")
	}
	
	if event.Status != "published" {
		return fmt.Errorf("event is not published")
	}

	var user domain.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	var count int64
	if err := s.DB.Model(&user).Where("id = ?", userID).
		Joins("JOIN user_events ON users.id = user_events.user_id").
		Where("user_events.event_id = ?", eventID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("error checking for existing event: %w", err)
	}

	if count > 0 {
			return fmt.Errorf("event already added to user")
	}

	user.AttendingEvents = append(user.AttendingEvents, event)
	err := s.DB.Model(&user).Association("AttendingEvents").Append(&event)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *MySQLAdapter) GetUsersEvents(userID string, status string) ([]domain.Event, error) {
	var user domain.User
	query := s.DB.Preload("AttendingEvents")
	
	if status == "active" {
		query = query.Preload("AttendingEvents", "date > ?", time.Now())
	} else if status == "completed" {
		query = query.Preload("AttendingEvents", "date <= ?", time.Now())
	}
	
	result := query.First(&user, "id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error fetching user: %w", result.Error)
	}
	
	return user.AttendingEvents, nil
}


