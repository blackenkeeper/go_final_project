package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blackenkeeper/go_final_project/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your_secret_key")

// Функция для создания JWT-токена
func createJWTToken() (string, error) {
	log.Println("Создание JWT-токена")
	token := jwt.New(jwt.SigningMethodHS256)

	log.Println("Подписание JWT-токена секретным ключём")
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Ошибка подписания токена:", err)
		return "", err
	}

	return tokenString, nil
}

// Функция для валидации JWT-токена
func validateJWTToken(tokenString string) (*jwt.Token, error) {
	log.Println("Валидация полученного токена")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			log.Println("Неожиданный метод подписи токена:", err)
			return nil, err
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		log.Println("Токен не валиден")
		return nil, err
	}

	return token, nil
}

// Middleware для аутентификации по JWT-токену.
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запуск middleware для аутентификации по JWT-токену")
		log.Println("Попытка получения токена из кук")
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println("Ошибка получения куки:", err)
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		jwtToken := cookie.Value
		_, err = validateJWTToken(jwtToken)
		if err != nil {
			log.Println("Токен не валиден:", err)
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}

// LoginHandler для обработки пути /task/signin и выдачи JWT-токена.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика LoginHandler для пути api/task/signin")

	// Пароль должен передаваться в JSON формате
	var creds struct {
		Password string `json:"password"`
	}

	log.Println("Попытка десериализации JSON и получения значения по ключу password из тела страницы")
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println("Ошибка получения данных из JSON:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Println("Попытка получения пароля из переменной окружения")
	pass := os.Getenv("TODO_PASSWORD")
	if creds.Password != pass {
		log.Println("Пароль из JSON не равен паролю в переменной окружения")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := createJWTToken()
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	log.Println("Сериализация токена в JSON для записи в тело страницы с ключём token")
	jsonAnswer, err := json.Marshal(&map[string]string{"token": token})
	if err != nil {
		log.Println("Ошибка сериализации токена в JSON:", err)
		ErrorsHandler(w, err, models.AnswerHandler{})
		return
	}

	log.Println("Обработчик LoginHandler отработал без ошибок, пользователю удалось авторизоваться")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(jsonAnswer); err != nil {
		ErrorsHandler(w, err, models.AnswerHandler{})
	}
}
