package handlers

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"net/http"
	"todolist/app/config"
	"todolist/app/internal/rules"
)

func DoneTask(w http.ResponseWriter, r *http.Request) {
	// проверяем, что метод запроса - Post
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	// получаем id задачи
	idTask := chi.URLParam(r, "id")
	if idTask == "" {
		idTask = r.URL.Query().Get("id")
	}

	// если id не указан, возвращаем ошибку
	if idTask == "" {
		http.Error(w, `{"error": "не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	var task Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := config.DB.QueryRow(query, idTask)
	if err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error": "ошибка получения задачи из БД"}`, http.StatusInternalServerError)
		}
		return
	}

	// проверяем, является ли задача периодической
	if task.Repeat == "" {
		// если задача одноразовая, удаляем её из базы
		deleteQuery := "DELETE FROM scheduler WHERE id = ?"
		_, err := config.DB.Exec(deleteQuery, idTask)
		if err != nil {
			http.Error(w, `{"error": "ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
	} else {
		// если задача периодическая, рассчитываем следующую дату выполнения
		nextDate, err := rules.NextDate(task.Date, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error": "ошибка расчёта следующей даты"}`, http.StatusInternalServerError)
			return
		}

		// обновляем дату задачи в базе данных
		updateQuery := "UPDATE scheduler SET date = ? WHERE id = ?"
		_, err = config.DB.Exec(updateQuery, nextDate, idTask)
		if err != nil {
			http.Error(w, `{"error": "ошибка обновления задачи"}`, http.StatusInternalServerError)
			return
		}
	}

	// возвращаем успешный пустой JSON
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}
