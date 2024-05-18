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
			"date char(8) not null default '' , ",
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
