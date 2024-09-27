package configs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
  _, err := gorm.Open(mysql.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }
	fmt.Println("🍀 MYSQL Connected")
}

