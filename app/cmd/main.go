package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"todolist/app/config"
	"todolist/app/internal/handlers"
)

var (
	defaultPort = "7540"
	webDir      = "./web"
)

func main() {
	config.LoadEnv()       // загружаем переменные окружения
	config.MakeDB()        // запуск БД
	defer config.CloseDB() // закрываем бд

	// проверка порта
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// запускаем роутер + логгер для запросов + восстановитель паники, чтобы не падал сервер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/nextdate", handlers.NextDate)
	r.Post("/api/task", handlers.PostTask)
	r.Get("/api/tasks", handlers.GetTasks)
	r.Get("/api/task", handlers.GetTask)
	r.Put("/api/task", handlers.PutTask)

	// добавляем обработчик файлов
	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// старт сервера
	log.Printf("Сервер запущен на http://localhost:%s/", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
