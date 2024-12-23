package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"todolist/app/config"
)

const (
	defaultPort = "7540"
	webDir      = "./web"
)

func main() {
	config.LoadEnv()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// запускаем роутер + логгер для запросов + восстановитель паники, чтобы не падал сервер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// добавляем обработчик файлов
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// старт сервера
	log.Printf("Сервер запущен на http://localhost:%s/", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
