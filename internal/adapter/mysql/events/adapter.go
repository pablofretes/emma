package adapters

import (
	events "emma/internal/core/domain/events"

	"gorm.io/gorm"
)

type MySQLAdapter struct {
	DB *gorm.DB
}

func NewMySQLAdapter(db *gorm.DB) (*MySQLAdapter) {
	return &MySQLAdapter{DB: db}
}

func (s *MySQLAdapter) CreateEvent(event events.Event) (events.Event, error) {
	result := s.DB.Create(&event)
	if result.Error != nil {
		return events.Event{}, result.Error
	}
	return event, nil
}

func (s *MySQLAdapter) GetEvent(id string) (events.Event, error) {
	var event events.Event
	result := s.DB.First(&event, "id = ?", id)
	if result.Error != nil {
		return events.Event{}, result.Error
	}
	return event, nil
}

func (s *MySQLAdapter) UpdateEvent(id string, event events.Event) (events.Event, error) {
	result := s.DB.Model(&events.Event{}).Where("id = ?", id).Updates(event)
	if result.Error != nil {
		return events.Event{}, result.Error
	}
	return s.GetEvent(id)
}

func (s *MySQLAdapter) DeleteEvent(id string) error {
	result := s.DB.Delete(&events.Event{}, "id = ?", id)
	return result.Error
}

func (s *MySQLAdapter) GetAllEvents() ([]events.Event, error) {
	var events []events.Event
	result := s.DB.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}
