package models

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/blackenkeeper/go_final_project/internal/config"
	"github.com/blackenkeeper/go_final_project/internal/repeater"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Response struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
	Tasks []Task `json:"tasks,omitempty"`
}

// Validate проверяет данные из задачи на валидность. Вторым аргументом принимает ссылку
// на объект, дата которого модифицируется при заданном значении task.Repeat функцией repeater.NextDate,
// если дата меньше сегодняшней. Иначе дата меняется на сегодняшнюю, если указанная в задаче дата раньше.
func (task *Task) Validate() (bool, error) {
	log.Debug("Запуск функции isAGoodTaskChecker")
	now := time.Now().Format(config.DateFormat)

	if task.Title == "" {
		log.Warn("Заголовок задачи не указан")
		return false, errors.New("не указан заголовок")
	}

	if task.Date == "" {
		task.Date = now
	}
	_, err := time.Parse(config.DateFormat, task.Date)
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
