package handlers

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// секретный ключ для подписи токенов, получаем его из переменной окружения
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

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
