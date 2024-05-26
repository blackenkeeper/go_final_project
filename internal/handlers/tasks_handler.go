package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/blackenkeeper/go_final_project/internal/database"
	"github.com/blackenkeeper/go_final_project/internal/models"
)

// Обработчик для пути /api/tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика TasksHandler для пути /api/tasks")
	w.Header().Set("Content-Type", "application/json")
	tasks := []models.Task{}

	var (
		answer      models.AnswerHandler
		selectQuery string
		rows        *sql.Rows
	)

	log.Println("Попытка подключения к базе данных по пути:", database.GetDbFile())
	db, err := sql.Open("sqlite", database.GetDbFile())
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	defer db.Close()

	searchParam := r.URL.Query().Get("search")
	limitParam := r.URL.Query().Get("limit")

	if limitParam == "" {
		limitParam = "50"
	}
	log.Printf("Получениe параметров из запроса: search - %s, limit - %s\n",
		searchParam, limitParam)

	log.Println("Попытка получения задач из базы данных с помощью запроса ниже")
	dateParam, err := taskDateParsing(searchParam)
	if err == nil {
		dateString := dateParam.Format("20060102")
		selectQuery = "select * from scheduler where date = ? limit ?;"
		log.Printf("Исполнение запроса к базе данных: \"%s\" c параметрами - %s, %s\n", selectQuery,
			searchParam, limitParam)
		rows, err = db.Query(selectQuery, dateString, limitParam)
	} else if searchParam == "" {
		selectQuery = "select * from scheduler order by date limit ?;"
		log.Printf("Исполнение запроса к базе данных: \"%s\" c параметром %s,\n", selectQuery, limitParam)
		rows, err = db.Query(selectQuery, limitParam)
	} else {
		selectQuery = "select * from scheduler where title like ? or comment like ? order by date limit ?;"
		searchParam = sqlLikeModder(searchParam)
		log.Printf("Исполнение запроса к базе данных: \"%s\" c параметрами - %s, %s, %s\n", selectQuery,
			searchParam, searchParam, limitParam)
		rows, err = db.Query(selectQuery, searchParam, searchParam, limitParam)
	}

	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		log.Println("Попытка чтения данных из возвращённой ссылки на *sql.Rows в пустой объект задачи")
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			ErrorsHandler(w, err, answer)
			return
		}
		log.Println("Пустой объект задачи после заполнения данными из *sql.Row:", task)
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Попытка сериализации полученных значений массива tasks и запись в тело ответа")
	tasksJson, err := json.Marshal(&map[string][]models.Task{"tasks": tasks})
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Обработчик TasksHandler отработал без ошибок и вернул список задач по запросу к БД")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(tasksJson); err != nil {
		ErrorsHandler(w, err, answer)
	}
}
