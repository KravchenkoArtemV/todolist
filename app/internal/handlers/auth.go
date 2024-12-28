package handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

// секретный ключ для подписи токенов, получаем его из переменной окружения
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

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

	// получаем пароль из переменной окружения
	passwordFromEnv := os.Getenv("TODO_PASSWORD")

	// проверяем: установлен ли пароль на сервере
	if passwordFromEnv == "" {
		// если переменная окружения не содержит пароль, возвращаем ошибку
		http.Error(w, `{"error":"пароль не установлен на сервере"}`, http.StatusInternalServerError)
		return
	}

	// проверяем введённый пользователем пароль
	if passwordFromEnv != enterPassword.Password {
		// если пароли не совпадают, возвращаем ошибку авторизации
		http.Error(w, `{"error":"некорректный пароль"}`, http.StatusUnauthorized)
		return
	}

	// создаём токен с использованием алгоритма подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password":   passwordFromEnv,                      // сохраняем текущий пароль как часть полезной нагрузки токена
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

// прослойка для проверки авторизации через токен
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// получаем пароль из переменной окружения
		envPassword := os.Getenv("TODO_PASSWORD")
		if envPassword == "" {
			// если пароль не задан в переменной окружения, пропускаем запрос без проверки
			next.ServeHTTP(w, r)
			return
		}

		// извлекаем токен из куки
		cookie, err := r.Cookie("token")
		if err != nil {
			// если кука с токеном отсутствует, возвращаем ошибку авторизации
			http.Error(w, `{"error":"необходима аутентификация"}`, http.StatusUnauthorized)
			return
		}

		// парсим токен из куки
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			// проверяем, что метод подписи токена - HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// если метод подписи некорректный, возвращаем ошибку подписи
				return nil, jwt.ErrSignatureInvalid
			}
			// возвращаем секретный ключ для проверки подписи
			return jwtSecret, nil
		})

		// проверяем валидность токена
		if err != nil || !token.Valid {
			// если токен недействителен, возвращаем ошибку авторизации
			http.Error(w, `{"error":"токен недействителен"}`, http.StatusUnauthorized)
			return
		}

		// извлекаем claims (полезную нагрузку) из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["password"] != envPassword {
			// если claims не совпадают с текущим паролем, токен считается недействительным
			http.Error(w, `{"error":"токен недействителен"}`, http.StatusUnauthorized)
			return
		}

		// если проверка токена успешна, передаём управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
