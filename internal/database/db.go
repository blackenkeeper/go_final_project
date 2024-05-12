package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func SetupDB() {
	var install bool
	_, err := os.Stat(GetDbFile())
	if err != nil {
		install = true
	}

	database, err := sql.Open("sqlite", GetDbFile())
	if err != nil {
		log.Fatalf("There's a problem with open db-file: %s\n", err)
	}
	defer database.Close()

	if install {
		createTable := fmt.Sprint("create table if not exists 'scheduler'(",
			"id integer primary key, ",
			"date varchar(8) not null default '' check(length(date) = 8), ",
			"title varchar(64) not null default '', ",
			"comment varchar(256) not null default '', ",
			"repeat varchar(128) not null default ''",
			");")

		_, err = database.Exec(createTable)
		if err != nil {
			log.Fatalf("Cannot create a table with db.Exec(): %s\n", err)
		}
		indexById := "create index id_index on scheduler (id);"
		_, err = database.Exec(indexById)
		if err != nil {
			log.Fatalf("Cannot create an index for table scheduler: %s\n", err)
		}
	}
}

func GetDbFile() string {
	envFilePath := os.Getenv("TODO_DBFILE")

	if envFilePath == "" {
		envFilePath = "../scheduler.db"
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	dbFilePath := filepath.Join(path, envFilePath)

	return dbFilePath
}
