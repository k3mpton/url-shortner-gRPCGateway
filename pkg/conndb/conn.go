package conndb

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	getenvfield "github.com/k3mpton/shortner-project/pkg/GetEnvField"
)

// Conn устанавливает соединение с базой данных PostgreSQL
// Использует драйвер pgx и получает строку подключения из переменной окружения CONN_DB
func Conn() *sql.DB {
	db, err := sql.Open("pgx", getenvfield.Get("CONN_DB"))
	if err != nil {
		log.Fatalf("не удалось обратиться по ссылке к бд: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v ", err)
	}

	return db
}
