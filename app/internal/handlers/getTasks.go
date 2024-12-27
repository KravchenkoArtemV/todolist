package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todolist/app/config"
	"todolist/app/internal/rules"
)

// получение списка задач
func GetTasks(w http.ResponseWriter, r *http.Request) {
	// проверяем, что метод ГЕТ
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}
	// заголовок
	w.Header().Set("Content-Type", "application/json")

	// получаем из строки запроса
	search := r.URL.Query().Get("search")

	// задаем переменные - лимит вывода, запрос, слайс аргументов к нему
	var limit = 50
	var query string
	var args []interface{}

	// если параметр search не указан, возвращаем все задачи
	if search == "" {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?"
		args = append(args, limit)
	} else {
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			// дата
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?"
			args = append(args, parsedDate.Format(rules.FormatTime), limit)
		} else {
			// подстрока
			searchPattern := "%" + search + "%"
			query = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?"
			args = append(args, searchPattern, searchPattern, limit)
		}
	}

	// выполняем запрос к БД
	rows, err := config.DB.Query(query, args...)
	if err != nil {
		http.Error(w, `{"error": "ошибка выполнения запроса к БД"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// слайс хранилище задач
	var tasks []Task
	for rows.Next() {
		var task Task
		// считываем данные из результата запроса в структуру
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			http.Error(w, `{"error": "ошибка сканирования БД"}`, http.StatusInternalServerError)
			return
		}
		// аппендим в слайс
		tasks = append(tasks, task)
	}
	// проверка на ошибки
	if err := rows.Err(); err != nil {
		http.Error(w, `{"error": "ошибка обработки БД"}`, http.StatusInternalServerError)
		return
	}
	// создаем слайс мап
	var tasksResp []map[string]string
	for _, task := range tasks {
		tasksResp = append(tasksResp, map[string]string{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		})
	}

	// возвращаем слайс если пустой ответ
	if tasksResp == nil {
		tasksResp = []map[string]string{}
	}

	// формируем ответ в JSON
	response := map[string]interface{}{"tasks": tasksResp}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, `{"error":"Ошибка формирования JSON-ответа"}`, http.StatusInternalServerError)
		return
	}
}
