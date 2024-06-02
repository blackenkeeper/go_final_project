package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/models"
	"github.com/blackenkeeper/go_final_project/internal/repeater"
	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func NewDB() (*Storage, error) {
	s := &Storage{}
	err := s.SetupDB()
	if err != nil {
		log.Fatal(err)
	}

	return s, err
}

// Настройка и запуск базы данных
func (s *Storage) SetupDB() error {
	var install bool
	_, err := os.Stat(GetDbFile())
	if err != nil {
		install = true
	}

	database, err := sql.Open("sqlite", GetDbFile())
	if err != nil {
		log.WithError(err).Error("Проблема с открытием файла базы данных")
		return err
	}

	if install {
		createTable := fmt.Sprint("CREATE TABLE IF NOT EXISTS 'scheduler'(",
			"id INTEGER PRIMARY KEY, ",
			"date CHAR(8) NOT NULL DEFAULT '' , ",
			"title varchar(64) NOT NULL DEFAULT '', ",
			"comment varchar(256) NOT NULL DEFAULT '', ",
			"repeat varchar(128) NOT NULL DEFAULT ''",
			");")

		_, err = database.Exec(createTable)
		if err != nil {
			log.WithError(err).Error("Невозможно создать базу данных вызовом db.Exec()")
			return err
		}
		indexById := "CREATE INDEX id_index ON scheduler (id);"
		_, err = database.Exec(indexById)
		if err != nil {
			log.WithError(err).Error("Ошибка создания индекса для таблицы scheduler")
			return err
		}
	}

	s.db = database
	return nil
}

func (s *Storage) GetTasks(searchParam, limitParam string) ([]models.Task, error) {
	db := s.db
	tasks := []models.Task{}

	var (
		selectQuery string
		rows        *sql.Rows
	)

	dateParam, err := taskDateParsing(searchParam)
	if err == nil {
		dateString := dateParam.Format("20060102")
		selectQuery = "SELECT * FROM scheduler WHERE date = ? LIMIT ?;"
		rows, err = db.Query(selectQuery, dateString, limitParam)
	} else if searchParam == "" {
		selectQuery = "SELECT * FROM scheduler ORDER BY date LIMIT ?;"
		rows, err = db.Query(selectQuery, limitParam)
	} else {
		selectQuery = "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?;"
		searchParam = sqlLikeModder(searchParam)
		rows, err = db.Query(selectQuery, searchParam, searchParam, limitParam)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, err
}

func (s *Storage) AddTask(task models.Task) (int, error) {
	if goodTask, err := isAGoodTaskChecker(&task); !goodTask || err != nil {
		return 0, errors.New("задача не соответствует заданному шаблону")
	}

	insertQuery := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?);"
	res, err := s.db.Exec(insertQuery, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (s *Storage) DeleteTask(id string) error {
	if _, err := strconv.Atoi(id); err != nil {
		return err
	}

	task, err := findTaskById(s, id)
	if err != nil {
		return err
	}

	deleteQuery := "DELETE FROM scheduler WHERE id = ?;"
	_, err = s.db.Exec(deleteQuery, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateTask(task models.Task) error {
	if goodTask, err := isAGoodTaskChecker(&task); !goodTask || err != nil {
		return errors.New("задача не соответствует заданному шаблону")
	}

	if _, err := findTaskById(s, task.ID); err != nil {
		return err
	}

	updateQuery := "UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?;"
	_, err := s.db.Exec(updateQuery, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) FindById(id string) (models.Task, error) {
	var task models.Task
	if _, err := strconv.Atoi(id); err != nil {
		return task, err
	}

	task, err := findTaskById(s, id)
	if err != nil {
		return task, err
	}

	return task, err
}

func (s *Storage) TaskDone(id string) (models.Task, error) {
	var task models.Task
	if _, err := strconv.Atoi(id); err != nil {
		return task, err
	}

	task, err := findTaskById(s, id)
	if err != nil {
		return task, err
	}

	if task.Repeat != "" {
		lastTaskDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			return task, err
		}

		task.Date, err = repeater.NextDate(lastTaskDate, task.Date, task.Repeat)
		if err != nil {
			return task, err
		}
	}

	return task, err
}

func (s *Storage) CloseDB() {
	s.db.Close()
}

// GetDbFile возвращает путь к базе данных. Расчитывается исходя из текущей рабочей директории,
// создаёт файл в ней, если его нет. Название файла можно задавать в переменной окружения TODO_DBFILE.
func GetDbFile() string {
	envFile := os.Getenv("TODO_DBFILE")

	if envFile == "" {
		envFile = "scheduler.db"
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	dbFilePath := filepath.Join(path, envFile)

	return dbFilePath
}
