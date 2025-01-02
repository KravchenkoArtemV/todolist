package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"todolist/app/config"

	"github.com/go-chi/chi/v5"
)

// обработчик для получения задачи по id
func GetTask(w http.ResponseWriter, r *http.Request) {
	// проверяем, что метод запроса - GET
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	// получаем id задачи из параметра пути или строки запроса
	idTask := chi.URLParam(r, "id")
	if idTask == "" {
		idTask = r.URL.Query().Get("id")
	}

	// если id не указан, возвращаем ошибку
	if idTask == "" {
		http.Error(w, `{"error": "не указан идентификатор"}`, http.StatusBadRequest)
		return
	}

	// SQL-запрос для получения задачи
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := config.DB.QueryRow(query, idTask)

	// структура для хранения результата
	var task Task
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// если задача не найдена, возвращаем ошибку
			http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		} else {
			// если произошла другая ошибка, возвращаем её
			http.Error(w, `{"error": "ошибка получения задачи из БД"}`, http.StatusInternalServerError)
		}
		return
	}

	// возвращаем найденную задачу в формате JSON
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, `{"error":"ошибка формирования JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
}
