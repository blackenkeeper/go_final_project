package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/blackenkeeper/go_final_project/internal/models"
)

// Функция для обработки ошибок и их записи в тело страницы.
// В третий параметр передаётся структура для ответов в формате JSON
func ErrorsHandler(w http.ResponseWriter, err error, answer models.AnswerHandler) {
	log.Println("Запуск обработчика ErrorsHandler для пути обработки полученных ошибок")
	log.Println("Поймана ошибка:", err)
	answer.Error = err.Error()
	bodyPage, _ := json.Marshal(answer)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(bodyPage)
	log.Println("Обработчик ErrorsHandler отработал успешно")
}
