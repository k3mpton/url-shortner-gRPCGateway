package getenvfield

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Get(key string) string {
	// Пытаемся загрузить .env файл, но не падаем если его нет
	fmt.Println("===========")
	fmt.Println(os.Getenv("DATABASE_URL"))
	if err := godotenv.Load(); err != nil {
		// Логируем предупреждение, но не падаем
		log.Printf("Предупреждение: .env файл не найден, используем переменные окружения: %v", err)
	}

	val := os.Getenv(key)
	fmt.Println("valll: ", val)
	if val == "" {
		log.Fatalf("Не удалось найти переменную окружения: %v", key)
	}

	return val
}
