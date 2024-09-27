package adapters

import (
	users "emma/internal/core/domain/users"

	"gorm.io/gorm"
)


type MySQLAdapter struct {
	DB *gorm.DB
}

func NewMySQLAdapter(db *gorm.DB) (*MySQLAdapter) {
	return &MySQLAdapter{DB: db}
}

func (s *MySQLAdapter) CreateUser(user users.User) error {
	result := s.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *MySQLAdapter) GetUserByUsername(username string) (users.User, error) {
	var user users.User
	result := s.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		return users.User{}, result.Error
	}
	return user, nil
}
