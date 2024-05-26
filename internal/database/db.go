package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/blackenkeeper/go_final_project/internal/utils"
	_ "modernc.org/sqlite"
)

// Настройка и запуск базы данных
func SetupDB() {
	var install bool
	_, err := os.Stat(GetDbFile())
	if err != nil {
		install = true
	}

	database, err := sql.Open("sqlite", GetDbFile())
	if err != nil {
		log.Fatalf("Проблема с открытием файла базы данных: %s\n", err)
	}
	defer database.Close()

	if install {
		createTable := fmt.Sprint("create table if not exists 'scheduler'(",
			"id integer primary key, ",
			"date char(8) not null default '' , ",
			"title varchar(64) not null default '', ",
			"comment varchar(256) not null default '', ",
			"repeat varchar(128) not null default ''",
			");")

		_, err = database.Exec(createTable)
		if err != nil {
			log.Fatalf("Невозможно создать базу данных вызовом db.Exec(): %s\n", err)
		}
		indexById := "create index id_index on scheduler (id);"
		_, err = database.Exec(indexById)
		if err != nil {
			log.Fatalf("Ошибка создания индекса для таблицы scheduler: %s\n", err)
		}
	}
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

	path = utils.CmdPathChecker(path)
	dbFilePath := filepath.Join(path, envFile)

	return dbFilePath
}
