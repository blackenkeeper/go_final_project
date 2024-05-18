package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/database"
	"github.com/blackenkeeper/go_final_project/internal/models"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tasks := []models.Task{}

	var (
		answer      models.AnswerHandler
		selectQuery string
		rows        *sql.Rows
	)

	db, err := sql.Open("sqlite", database.GetDbFile())
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	defer db.Close()

	searchParam := r.URL.Query().Get("search")
	limitParam := r.URL.Query().Get("limit")

	if limitParam == "" {
		limitParam = "10"
	}
	log.Println("Search parameter:", searchParam)

	dateParam, err := dateParsing(searchParam)
	if err == nil {
		dateString := dateParam.Format("20060102")
		selectQuery = "select * from scheduler where date = ? limit ?;"
		log.Println("Executing query:", selectQuery, "with params:", dateString, limitParam)
		rows, err = db.Query(selectQuery, dateString, limitParam)
	} else if searchParam == "" {
		selectQuery = "select * from scheduler order by date limit ?;"
		log.Println("Executing query:", selectQuery, "with params:", limitParam)
		rows, err = db.Query(selectQuery, limitParam)
	} else {
		selectQuery = "select * from scheduler where title like ? or comment like ? order by date limit ?;"
		searchParam = forLikeModder(searchParam)
		log.Println("Executing query:", selectQuery, "with params:", searchParam, searchParam, limitParam)
		rows, err = db.Query(selectQuery, searchParam, searchParam, limitParam)
	}

	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			ErrorsHandler(w, err, answer)
			return
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	tasksJson, err := json.Marshal(&map[string][]models.Task{"tasks": tasks})
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(tasksJson); err != nil {
		ErrorsHandler(w, err, answer)
	}
}

func dateParsing(date string) (time.Time, error) {
	return time.Parse(`02.01.2006`, date)
}

func forLikeModder(s string) string {
	return "%" + s + "%"
}
