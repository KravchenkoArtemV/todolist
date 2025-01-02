package config

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var PasswordFromEnv string

// загрузка переменных окружения
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}
}

// создание бд
func MakeDB() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	// загружаем переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = filepath.Join(filepath.Dir(appPath), "database", "scheduler.db")
	}

	log.Printf("Путь к БД: %s", dbFile)
	// проверка существования базы данных
	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// если база данных новая, создаём таблицы и индексы
	if install {
		err = createTable(DB)
		if err != nil {
			log.Fatalf("Ошибка создания таблицы: %v", err)
		}
		log.Println("БД успешно создана")
	}
}

// создаем таблицы и индексы в базе данных
func createTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS scheduler (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date TEXT NOT NULL,
        title TEXT NOT NULL,
        comment TEXT,
        repeat TEXT CHECK (LENGTH(repeat) <= 128)
    );

    CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
    `
	_, err := db.Exec(query)
	return err
}

// разрываем соединение с бд
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Ошибка закрытия БД: %v", err)
		}
		log.Println("Соединение с БД закрыто")
	}
}
