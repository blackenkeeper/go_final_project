package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/database"
	"github.com/blackenkeeper/go_final_project/internal/models"
	"github.com/blackenkeeper/go_final_project/internal/repeater"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {

	requestType := r.Method

	switch requestType {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		task models.Task
		buf  bytes.Buffer
	)
	answer := models.HandlerAnswer{}

	db, err := sql.Open("sqlite", database.GetDbFile())
	if err != nil {
		log.Fatalf("Не установлено соединение с базой данных по причине: %s", err)
	}
	defer db.Close()

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		errorsHandler(w, err, answer)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		errorsHandler(w, err, answer)
		return
	}
	if !taskChecker(w, answer, &task) {
		return
	}
	insertQuery := "insert into scheduler (date, title, comment, repeat) values (?, ?, ?, ?);"
	res, err := db.Exec(insertQuery, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		errorsHandler(w, err, answer)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		errorsHandler(w, err, answer)
		return
	}
	answer.ID = int(id)

	bodyPage, err := json.Marshal(answer)
	if err != nil {
		errorsHandler(w, err, answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bodyPage)

}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
}

func errorsHandler(w http.ResponseWriter, err error, answer models.HandlerAnswer) {
	answer.Error = err.Error()
	bodyPage, _ := json.Marshal(answer)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(bodyPage)
}

func taskChecker(w http.ResponseWriter, answer models.HandlerAnswer, task *models.Task) bool {
	now := time.Now().Format("20060102")

	if task.Title == "" {
		errorsHandler(w, errors.New("не указан заголовок"), answer)
		return false
	}
	if task.Date == "" {
		task.Date = now
	}
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		errorsHandler(w, err, answer)
		return false
	}

	if task.Repeat != "" && task.Date != now {
		task.Date, err = repeater.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			errorsHandler(w, err, answer)
			return false
		}
	} else {
		task.Date = now
	}

	return true
}
