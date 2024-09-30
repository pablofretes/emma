package dbs

import (
	configs "emma/configs"
	"emma/internal/domain"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dbName := configs.GetConfig().MYSQL_DB
	dbUser := configs.GetConfig().MYSQL_USER
	dbPassword := configs.GetConfig().MYSQL_PASSWORD
	dbTcp := configs.GetConfig().MYSQL_TCP
	dbDriver := dbUser+":"+dbPassword+dbTcp+"/"+dbName+
	"?charset=utf8&parseTime=True"
	db, err := gorm.Open(mysql.Open(dbDriver), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database", err)
	}
	fmt.Println("🍀 MYSQL Connected")

    err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&domain.User{}, &domain.Event{})
    if err != nil {
        log.Fatal("Error migrating database:", err)
    }

	fmt.Println("🍀 Database and tables created/updated successfully")

	return db
}
