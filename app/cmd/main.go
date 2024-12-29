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
	config.LoadEnv()
	config.MakeDB()
	defer config.CloseDB()

	// проверка порта
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// запускаем роутер + логгер для запросов + восстановитель паники, чтобы не падал сервер
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/*", http.FileServer(http.Dir(webDir)))

	// обработчик для авторизации
	r.Post("/api/signin", handlers.SignIn)

	// защищенные маршруты будут прогоняться через токен
	r.Route("/api", func(r chi.Router) {
		// добавляем middleware
		r.Use(handlers.Auth)

		r.Get("/nextdate", handlers.NextDate)
		r.Post("/task", handlers.PostTask)
		r.Get("/tasks", handlers.GetTasks)
		r.Get("/task", handlers.GetTask)
		r.Put("/task", handlers.PutTask)
		r.Post("/task/done", handlers.DoneTask)
		r.Delete("/task", handlers.RemoveTask)
	})

	// старт сервера
	log.Printf("Сервер запущен на http://localhost:%s/", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
