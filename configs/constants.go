package configs

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT            string `example:":8080"`
	MYSQL_DB        string `example:"myapp_db"`
	MYSQL_TCP       string `example:"@tcp(127.0.0.1:3306)"`
	MYSQL_USER      string `example:"db_user"`
	MYSQL_PASSWORD  string `example:"secret_password"`
	MYSQL_HOST      string `example:"localhost"`
	MYSQL_PORT      string `example:"3306"`
	JWT_SECRET      string `example:"my_secret_key"`
	JWT_EXPIRY_TIME string `example:"24h"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		instance = &Config{
			PORT: os.Getenv("PORT"),
			MYSQL_DB: os.Getenv("MYSQL_DB"),
			MYSQL_TCP: os.Getenv("MYSQL_TCP"),
			MYSQL_USER: os.Getenv("MYSQL_USER"),
			MYSQL_PASSWORD: os.Getenv("MYSQL_PASSWORD"),
			MYSQL_HOST: os.Getenv("MYSQL_HOST"),
			MYSQL_PORT: os.Getenv("MYSQL_PORT"),
			JWT_SECRET: os.Getenv("JWT_SECRET"),
			JWT_EXPIRY_TIME: os.Getenv("JWT_EXPIRY_TIME"),
		}
	})
	return instance
}
