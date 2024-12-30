package handlers

// структура задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// структура ответа
type Response struct {
	Id    string `json:"id"`
	Error string `json:"error"`
}
