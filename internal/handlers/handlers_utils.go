package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/database"
	"github.com/blackenkeeper/go_final_project/internal/models"
	"github.com/blackenkeeper/go_final_project/internal/repeater"
)

// isAGoodTaskChecker проверяет данные из задачи на валидность. Вторым аргументом принимает ссылку
// на объект, дата которого модифицируется при заданном значении task.Repeat функцией repeater.NextDate,
// если дата меньше сегодняшней. Иначе дата меняется на сегодняшнюю, если указанная в задаче дата раньше.
func isAGoodTaskChecker(w http.ResponseWriter, task *models.Task) bool {
	log.Println("Запуск функции isAGoodTaskChecker")
	log.Printf("Проверка значний задачи %s на валидность и расчёт следующей даты повторения задачи\n", task)
	now := time.Now().Format("20060102")
	answer := models.AnswerHandler{}

	if task.Title == "" {
		ErrorsHandler(w, errors.New("не указан заголовок"), answer)
		return false
	}
	if task.Date == "" {
		task.Date = now
	}
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return false
	}

	if task.Repeat != "" && task.Date < now {
		task.Date, err = repeater.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			ErrorsHandler(w, err, answer)
			return false
		}
	} else if task.Date < now {
		task.Date = now
	}

	log.Println("Успешное окончание работы функции isAGoodTaskChecker")
	return true
}

// findTaskById служит для поиска задачи по заданному параметру id и её возвращения в случае
// успешного нахождения. Возвращает пустую задачу и ошибку в случае error != nil.
// Если параметр db == nil, то создаёт новое подключение к базе данных.
func findTaskById(db *sql.DB, w http.ResponseWriter, id string) (models.Task, error) {
	log.Println("Запуск функции findTaskById с id:", id)
	var (
		err    error
		answer models.AnswerHandler
	)

	if db == nil {
		log.Println("Попытка подключения к базе данных по пути:", database.GetDbFile())
		db, err = sql.Open("sqlite", database.GetDbFile())
		if err != nil {
			ErrorsHandler(w, err, answer)
			return models.Task{}, err
		}
	}

	selectQuery := "select * from scheduler where id = ?;"
	log.Printf("Исполнение запроса к базе данных: \"%s\" c параметром %s", selectQuery, id)
	taskRow := db.QueryRow(selectQuery, id)

	task := models.Task{}
	log.Println("Попытка чтения данных из возвращённой ссылки на *sql.Row в пустой объект задачи")
	err = taskRow.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		ErrorsHandler(w, err, answer)
		return models.Task{}, err
	}
	log.Println("Пустой объект задачи после заполнения данными из *sql.Row:", task)
	log.Println("Успешное окончание работы функции findTaskById")

	return task, err
}

// taskDateParsing используется для проверки соотвествия строки даты используемому шаблону.
func taskDateParsing(date string) (time.Time, error) {
	return time.Parse(`02.01.2006`, date)
}

// sqlLikeModder служит для модификации строки для исполнения sql-запроса с параметром like (для поиска по
// подстроке s).
func sqlLikeModder(s string) string {
	return "%" + s + "%"
}
