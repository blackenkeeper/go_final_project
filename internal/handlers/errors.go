package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/blackenkeeper/go_final_project/internal/models"
)

// Функция для обработки ошибок и их записи в тело страницы.
// В третий параметр передаётся структура для ответов в формате JSON
func (h *Handler) ErrorsHandler(w http.ResponseWriter, err error, answer models.Response) {
	log.Error("Поймана ошибка:", err)
	answer.Error = err.Error()
	bodyPage, _ := json.Marshal(answer)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(bodyPage)
}
