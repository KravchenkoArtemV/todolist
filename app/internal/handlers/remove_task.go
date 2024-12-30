package handlers

import (
	"net/http"
	"todolist/app/config"

	"github.com/go-chi/chi/v5"
)

// обработчик для удаления задачи
func RemoveTask(w http.ResponseWriter, r *http.Request) {
	// проверяем, что метод запроса - DELETE
	if r.Method != http.MethodDelete {
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
		http.Error(w, `{"error": "не указан идентификатор задачи"}`, http.StatusBadRequest)
		return
	}

	// выполняем запрос на удаление задачи из базы данных
	deleteQuery := "DELETE FROM scheduler WHERE id = ?"
	res, err := config.DB.Exec(deleteQuery, idTask)
	if err != nil {
		http.Error(w, `{"error": "ошибка удаления задачи"}`, http.StatusInternalServerError)
		return
	}

	// проверяем, была ли удалена хотя бы одна строка
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error": "ошибка проверки удаления задачи"}`, http.StatusInternalServerError)
		return
	}

	// если ни одна строка не была затронута, задача с таким id не найдена
	if rowsAffected == 0 {
		http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		return
	}

	// возвращаем успешный пустой JSON
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("{}"))
}
