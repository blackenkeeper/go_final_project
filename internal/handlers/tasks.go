package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/blackenkeeper/go_final_project/internal/models"
	log "github.com/sirupsen/logrus"
)

// Обработчик для пути /api/tasks
func (h *Handler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Запуск обработчика TasksHandler для пути /api/tasks")
	w.Header().Set("Content-Type", "application/json")

	answer := models.Response{}

	searchParam := r.URL.Query().Get("search")
	limitParam := r.URL.Query().Get("limit")
	if limitParam == "" {
		limitParam = "50"
	}

	tasks, err := h.Storage.GetTasks(searchParam, limitParam)
	if err != nil {
		log.WithError(err).Error("Ошибка при получении задач")
		h.ErrorsHandler(w, err, answer)
		return
	}

	tasksJson, err := json.Marshal(&map[string][]models.Task{"tasks": tasks})
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	log.Info("Обработчик TasksHandler отработал без ошибок и вернул список задач по запросу к БД")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(tasksJson); err != nil {
		h.ErrorsHandler(w, err, answer)
	}
}
