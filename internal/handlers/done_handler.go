package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/models"
	"github.com/blackenkeeper/go_final_project/internal/repeater"
)

// TaskDoneHandler для обработки пути /api/task/done
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика TaskDoneHandler для пути api/task/done")
	w.Header().Set("Content-Type", "application/json")

	var answer models.AnswerHandler

	log.Println("Попытка получения параметра id из запроса")
	taskId := r.URL.Query().Get("id")
	if _, err := strconv.Atoi(taskId); err != nil {
		log.Println("Ошибка конвертации в число параметра id:", err)
		ErrorsHandler(w, err, answer)
		return
	}

	task, err := findTaskById(nil, w, taskId)
	if err != nil {
		return
	}

	if task.Repeat != "" {
		log.Printf("Парсинг строки с датой из задачи %s: %s\n", task.Title, task.Date)
		lastTaskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			ErrorsHandler(w, err, answer)
			return
		}

		log.Println("Попытка расчёта следующей даты повторения задачи по правилу:", task.Repeat)
		task.Date, err = repeater.NextDate(lastTaskDate, task.Date, task.Repeat)
		if err != nil {
			ErrorsHandler(w, err, answer)
			return
		}
	}

	log.Println("Попытка сериализации задачи:", task)
	taskJson, err := json.Marshal(&task)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Проверка правила повторения:", task.Repeat)
	switch task.Repeat {
	case "":
		log.Println("Значение для повторения не задано, попытка удалить задачу вызовом deleteTaskHandler")
		newRequest, err := http.NewRequest(http.MethodDelete, r.URL.String(), bytes.NewBuffer(taskJson))
		if err != nil {
			ErrorsHandler(w, err, answer)
			return
		}

		deleteTaskHandler(w, newRequest)
	default:
		log.Println("Расчёт следующей даты задачи", task.Title, "по правилу", task.Repeat,
			"для передачи в обработчик updateTaskHandler и обновления даты задачи на", task.Date)
		newRequest, err := http.NewRequest(http.MethodPut, r.URL.String(), bytes.NewBuffer(taskJson))
		if err != nil {
			ErrorsHandler(w, err, answer)
			return
		}

		updateTaskHandler(w, newRequest)
	}

	log.Println("Обработчик TaskDoneHandler отработал без ошибок")

}
