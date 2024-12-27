package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todolist/app/config"
	"todolist/app/internal/rules"
)

// добавление задачи
func PostTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	// заголовок ответа
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// создаем объект
	var task Task
	// десериализуем json
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "ошибка десериализации"})
		return
	}
	// проверка наличия заполненного поля title
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "не указано title"})
		return
	}
	// если не указана дата, присваиваем сегодняшнюю
	nowDate := time.Now().Format(rules.FormatTime)
	if task.Date == "" || task.Date == "today" {
		task.Date = nowDate
	} else {
		// парсим дату
		parsedDate, err := time.Parse(rules.FormatTime, task.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Error: "неверный формат даты"})
			return
		}
		// если дата меньше сегодняшней
		if parsedDate.Before(time.Now()) {
			// правило не указано
			if task.Repeat == "" {
				task.Date = nowDate
			} else {
				// вычисляем следующую дату выполнения
				nextDate, err := rules.NextDate(nowDate, task.Date, task.Repeat)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(Response{Error: "ошибка в правиле повторения: "})
					return
				}

				if task.Date != nowDate {
					task.Date = nextDate
				}
			}
		}
	}
	// проверяем правило по умолчанию
	if task.Repeat != "" {
		_, err := rules.NextDate(nowDate, task.Date, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Некорректное правило повторения"}`, http.StatusBadRequest)
			return
		}
	}

	// добавляем задачу в базу данных
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"
	res, err := config.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "ошибка при добавлении задачи в БД"})
		return
	}

	// Получаем ID добавленной задачи
	id, err := res.LastInsertId()
	// делаем стрингу
	idResp := strconv.Itoa(int(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Error: "Ошибка при получении ID задачи"})
		return
	}

	// Отправляем успешный ответ с ID задачи
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": idResp})
}
