package config

import (
	"github.com/joho/godotenv"
	"log"
)

// загрузка переменных окружения
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
}
