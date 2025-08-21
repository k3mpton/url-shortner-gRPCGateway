package getenvfield

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Get(key string) string {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Не удалось спарсить .env файл: %v", err)
	}

	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Не удалось найти ключ: %v", key)
	}

	return val
}
