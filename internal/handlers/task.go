package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/blackenkeeper/go_final_project/internal/models"
)

// Планировщик для обработчиков пути /api/task. Для каждого типа запроса передаёт обработку
// нужному обработчику через switch по типу метода
func (h *Handler) TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTaskByIdHandler(w, r)
	case http.MethodPost:
		h.addTaskHandler(w, r)
	case http.MethodPut:
		h.updateTaskHandler(w, r)
	case http.MethodDelete:
		h.deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик POST-запроса к /api/task
func (h *Handler) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Запуск обработчика addTaskHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var (
		task   models.Task
		buf    bytes.Buffer
		answer models.AnswerHandler
	)

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}
	answer.ID, err = h.Storage.AddTask(task)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}
	bodyPage, err := json.Marshal(answer)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bodyPage); err != nil {
		h.ErrorsHandler(w, err, answer)
	}

	log.Info("Задача добавлена в базу данных")
}

// Обработчик DELETE-запроса к /api/task
func (h *Handler) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Запуск обработчика deleteTaskHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var answer models.AnswerHandler

	taskId := r.URL.Query().Get("id")
	err := h.Storage.DeleteTask(taskId)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{}")); err != nil {
		h.ErrorsHandler(w, err, answer)
	}

	log.Info("Задача удалена из базы данных")
}

// Обработчик PUT-запроса к /api/task
func (h *Handler) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Запуск обработчика updateTaskHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var (
		answer models.AnswerHandler
		buf    bytes.Buffer
		task   models.Task
	)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}
	err = h.Storage.UpdateTask(task)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{}")); err != nil {
		h.ErrorsHandler(w, err, answer)
	}

	log.Info("Задача обновлена в базе данных")
}

// Обработчик GET-запроса к /api/task. Необходим параметр id в запросе
func (h *Handler) getTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Запуск обработчика getTaskByIdHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var answer models.AnswerHandler

	taskId := r.URL.Query().Get("id")
	task, err := h.Storage.FindById(taskId)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	taskJson, err := json.Marshal(&task)
	if err != nil {
		h.ErrorsHandler(w, err, answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(taskJson); err != nil {
		h.ErrorsHandler(w, err, answer)
	}

	log.Info("Задача найдена по id")
}
