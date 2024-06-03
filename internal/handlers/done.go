package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/blackenkeeper/go_final_project/internal/models"
)

// TaskDoneHandler для обработки пути /api/task/done
func (h *Handler) TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Запуск обработчика TaskDoneHandler для пути api/task/done")
	w.Header().Set("Content-Type", "application/json")

	var answer models.Response

	taskId := r.URL.Query().Get("id")
	task, err := h.Storage.TaskDone(taskId)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	taskJson, err := json.Marshal(&task)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	switch task.Repeat {
	case "":
		newRequest, err := http.NewRequest(http.MethodDelete, r.URL.String(), bytes.NewBuffer(taskJson))
		if err != nil {
			h.ErrorsHandler(w, err, answer)
			return
		}

		h.deleteTaskHandler(w, newRequest)
	default:
		newRequest, err := http.NewRequest(http.MethodPut, r.URL.String(), bytes.NewBuffer(taskJson))
		if err != nil {
			h.ErrorsHandler(w, err, answer)
			return
		}

		h.updateTaskHandler(w, newRequest)
	}

	log.Info("Обработчик TaskDoneHandler отработал без ошибок")
}
