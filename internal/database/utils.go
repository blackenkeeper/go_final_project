package database

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/blackenkeeper/go_final_project/internal/models"
	"github.com/blackenkeeper/go_final_project/internal/repeater"
)

// isAGoodTaskChecker проверяет данные из задачи на валидность. Вторым аргументом принимает ссылку
// на объект, дата которого модифицируется при заданном значении task.Repeat функцией repeater.NextDate,
// если дата меньше сегодняшней. Иначе дата меняется на сегодняшнюю, если указанная в задаче дата раньше.
func isAGoodTaskChecker(task *models.Task) (bool, error) {
	log.Debug("Запуск функции isAGoodTaskChecker")
	now := time.Now().Format("20060102")

	if task.Title == "" {
		log.Warn("Заголовок задачи не указан")
		return false, errors.New("не указан заголовок")
	}
	if task.Date == "" {
		task.Date = now
	}
	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Warn("Некорректная дата задачи:", err)
		return false, err
	}

	if task.Repeat != "" && task.Date < now {
		task.Date, err = repeater.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Warn("Ошибка при расчете следующей даты повторения:", err)
			return false, err
		}
	} else if task.Date < now {
		task.Date = now
	}

	log.Debug("Успешное окончание работы функции isAGoodTaskChecker")
	return true, err
}

// findTaskById служит для поиска задачи по заданному параметру id и её возвращения в случае
// успешного нахождения. Возвращает пустую задачу и ошибку в случае error != nil.
func findTaskById(s *Storage, id string) (models.Task, error) {
	log.Debug("Запуск функции findTaskById с id:", id)
	var (
		err error
	)

	selectQuery := "select * from scheduler where id = ?;"
	log.Debug("Исполнение запроса к базе данных с параметром id:", id)
	taskRow := s.db.QueryRow(selectQuery, id)

	task := models.Task{}
	err = taskRow.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Error("Ошибка при чтении данных из базы данных:", err)
		return models.Task{}, err
	}

	log.Debug("Успешное окончание работы функции findTaskById")
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
