package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todolist/app/config"
	"todolist/app/internal/rules"
)

func PutTask(w http.ResponseWriter, r *http.Request) {
	// проверяем, что метод запроса - Put
	if r.Method != http.MethodPut {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var task Task
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		http.Error(w, `{"error":"Не указан ID"}`, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, `{"error":"Не указан заголовок"}`, http.StatusBadRequest)
		return
	}

	nowDate := time.Now().Format(rules.FormatTime)
	if task.Date == "" || task.Date == "today" {
		task.Date = nowDate
	} else {
		parsedDate, err := time.Parse(rules.FormatTime, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Некорректный формат даты"}`, http.StatusBadRequest)
			return
		}
		if parsedDate.Before(time.Now()) {
			if task.Repeat == "" {
				task.Repeat = nowDate
			} else {
				nextDate, err := rules.NextDate(nowDate, task.Date, task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
					return
				}
				if task.Date != nowDate {
					task.Date = nextDate
				}
			}
		}
	}
	// проверяем правило повторения
	if task.Repeat != "" {
		_, err := rules.NextDate(nowDate, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}
	// обновляем данные задачи в базе данных
	query := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?"
	res, err := config.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		http.Error(w, `{"error": "ошибка обновления задачи в базе данных"}`, http.StatusInternalServerError)
		return
	}

	// проверяем, была ли обновлена хотя бы одна строка
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, `{"error": "ошибка получения данных обновления"}`, http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, `{"error": "задача не найдена"}`, http.StatusNotFound)
		return
	}

	// возвращаем успешный пустой JSON
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
