package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todolist/app/config"

	"github.com/golang-jwt/jwt/v4"
)

// обработчик для входа пользователя (создание JWT токена)
func SignIn(w http.ResponseWriter, r *http.Request) {
	// проверяем, что метод запроса - POST
	if r.Method != http.MethodPost {
		// если метод не POST, возвращаем ошибку
		http.Error(w, `{"error": "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// создаем анонимную структуру для извлечения пароля из тела запроса
	var enterPassword struct {
		Password string `json:"password"` // поле, куда будет декодирован пароль из запроса
	}

	// декодируем тело запроса (JSON) в нашу структуру enterPassword
	if err := json.NewDecoder(r.Body).Decode(&enterPassword); err != nil {
		// если ошибка в декодировании, возвращаем ошибку клиенту
		http.Error(w, `{"error":"некорректный запрос"}`, http.StatusBadRequest)
		return
	}

	// проверяем: установлен ли пароль на сервере
	if config.PasswordFromEnv == "" {
		// если переменная окружения не содержит пароль, возвращаем ошибку
		http.Error(w, `{"error":"пароль не установлен на сервере"}`, http.StatusInternalServerError)
		return
	}

	// проверяем введённый пользователем пароль
	if config.PasswordFromEnv != enterPassword.Password {
		// если пароли не совпадают, возвращаем ошибку авторизации
		http.Error(w, `{"error":"некорректный пароль"}`, http.StatusUnauthorized)
		return
	}

	// создаём токен с использованием алгоритма подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password":   config.PasswordFromEnv,               // сохраняем текущий пароль как часть полезной нагрузки токена
		"expiration": time.Now().Add(8 * time.Hour).Unix(), // срок действия токена: 8 часов
	})

	// подписываем токен секретным ключом
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		// если произошла ошибка при создании токена, возвращаем ошибку
		http.Error(w, `{"error":"Ошибка генерации токена"}`, http.StatusInternalServerError)
		return
	}

	// устанавливаем заголовок ответа как JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// отправляем токен в теле ответа
	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}
