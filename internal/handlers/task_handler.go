package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/blackenkeeper/go_final_project/internal/database"
	"github.com/blackenkeeper/go_final_project/internal/models"
)

// Планировщик для обработчиков пути /api/task. Для каждого типа запроса передаёт обработку
// нужному обработчику через switch по типу метода
func TaskHandler(w http.ResponseWriter, r *http.Request) {

	requestType := r.Method

	switch requestType {
	case http.MethodGet:
		getTaskByIdHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

// Обработчик POST-запроса к /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика addTaskHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var (
		task models.Task
		buf  bytes.Buffer
	)
	answer := models.AnswerHandler{}

	log.Println("Чтение данных из тела страницы в буфер байт")
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Попытка десериализации считанных байт в объект структуры задачи")
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	if !isAGoodTaskChecker(w, &task) {
		return
	}

	log.Println("Попытка подключения к базе данных по пути:", database.GetDbFile())
	db, err := sql.Open("sqlite", database.GetDbFile())
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	defer db.Close()

	insertQuery := "insert into scheduler (date, title, comment, repeat) values (?, ?, ?, ?);"
	log.Printf("Исполнение запроса к базе данных: \"%s\" c параметрами - %s, %s, %s, %s\n", insertQuery,
		task.Date, task.Title, task.Comment, task.Repeat)

	res, err := db.Exec(insertQuery, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Получения последней добавленной строки в базу данных функцией LastInsertId()")
	id, err := res.LastInsertId()
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	answer.ID = int(id)

	log.Printf("Попытка сериализации полученного значения id(%d) и запись в тело ответа\n", id)
	bodyPage, err := json.Marshal(answer)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Обработчик addTaskHandler отработал без ошибок, задача добавлена в базу данных")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bodyPage); err != nil {
		ErrorsHandler(w, err, answer)
	}

}

// Обработчик DELETE-запроса к /api/task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика deleteTaskHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var answer models.AnswerHandler

	taskId := r.URL.Query().Get("id")
	log.Println("Получение параметра id из запроса:", taskId)
	if _, err := strconv.Atoi(taskId); err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Попытка подключения к базе данных по пути:", database.GetDbFile())
	db, err := sql.Open("sqlite", database.GetDbFile())
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	task, err := findTaskById(db, w, taskId)
	if err != nil {
		return
	}

	deleteQuery := "delete from scheduler where id = ?;"
	log.Printf("Попытка удалить задачу запросом: \"%s\" c параметром %s\n", deleteQuery, task.ID)
	_, err = db.Exec(deleteQuery, task.ID)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Обработчик deleteTaskHandler отработал без ошибок, задача удалена из базы данных")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("{}")); err != nil {
		ErrorsHandler(w, err, answer)
	}

}

// Обработчик PUT-запроса к /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика updateTaskHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var (
		answer models.AnswerHandler
		buf    bytes.Buffer
		task   models.Task
	)

	log.Println("Попытка подключения к базе данных по пути:", database.GetDbFile())
	db, err := sql.Open("sqlite", database.GetDbFile())
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Чтение данных из тела страницы в буфер байт")
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Попытка десериализации считанных байт в объект структуры задачи")
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	if !isAGoodTaskChecker(w, &task) {
		return
	}

	if _, err = findTaskById(db, w, task.ID); err != nil {
		return
	}

	updateQuery := "update scheduler set date = ?, title = ?, comment = ?, repeat = ? where id = ?;"
	log.Printf("Исполнение запроса к базе данных: \"%s\" c параметрами - %s, %s, %s, %s\n", updateQuery,
		task.Date, task.Title, task.Comment, task.Repeat)
	_, err = db.Exec(updateQuery, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}
	log.Println("Данные в БД были обновлены")

	log.Println("Обработчик updateTaskHandler отработал без ошибок")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{}")); err != nil {
		ErrorsHandler(w, err, answer)
	}
}

// Обработчик GET-запроса к /api/task. Необходим параметр id в запросе
func getTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика updateTaskHandler для пути /api/task")
	w.Header().Set("Content-Type", "application/json")

	var answer models.AnswerHandler

	taskId := r.URL.Query().Get("id")
	log.Println("Получение параметра id из запроса:", taskId)
	if _, err := strconv.Atoi(taskId); err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	task, err := findTaskById(nil, w, taskId)
	if err != nil {
		return
	}

	log.Printf("Попытка сериализации объекта task(%s) и запись в тело ответа\n", task)
	taskJson, err := json.Marshal(&task)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return
	}

	log.Println("Обработчик getTaskByIdHandler отработал без ошибок и нашёл задачу по заданному id")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(taskJson); err != nil {
		ErrorsHandler(w, err, answer)
	}
}
