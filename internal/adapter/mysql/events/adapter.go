package adapters

import (
	"emma/internal/domain"
	"strings"
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

func (s *MySQLAdapter) CreateEvent(event domain.Event) (domain.Event, error) {
	event.Attendees = []domain.User{}
	event.ID = uuid.New()
	event.Title = strings.ToLower(event.Title)
	event.ShortDescription = strings.ToLower(event.ShortDescription)
	event.LongDescription = strings.ToLower(event.LongDescription)
	event.Location = strings.ToLower(event.Location)
	event.Status = strings.ToLower(event.Status)
	result := s.DB.Create(&event)
	if result.Error != nil {
		return domain.Event{}, result.Error
	}
	return event, nil
}

func (s *MySQLAdapter) GetEvent(id string) (domain.Event, error) {
	var event domain.Event
	result := s.DB.First(&event, "id = ?", id)
	if result.Error != nil {
		return domain.Event{}, result.Error
	}

	return event, nil
}

func (s *MySQLAdapter) UpdateEvent(id string, event domain.Event) (domain.Event, error) {
	result := s.DB.Model(&domain.Event{}).Where("id = ?", id).Updates(event)
	if result.Error != nil {
		return domain.Event{}, result.Error
	}
	return s.GetEvent(id)
}

func (s *MySQLAdapter) DeleteEvent(id string) error {
	tx := s.DB.Begin()

	if err := tx.Table("user_events").Where("event_id = ?", id).Delete(nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&domain.Event{}, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *MySQLAdapter) GetAllEvents(status, title, date string) ([]domain.Event, error) {
	var events []domain.Event
	query := s.DB

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if date != "" {
		parsedDate, err := time.Parse("2006-01-02T15:04:05Z07:00", date)
		if err != nil {
			return nil, err
		}
		query = query.Where("date = ?", parsedDate)
	}

	result := query.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}
