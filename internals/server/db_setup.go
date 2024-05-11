package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func setupDB() {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "../scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("There's a problem with open db-file: %s\n", err)
	}
	defer db.Close()

	if install {
		createTable := fmt.Sprint("create table if not exists 'scheduler'(",
			"id integer primary key,",
			"date varchar(8) not null default '' check(length(date) = 8),",
			"title varchar(64) not null default '',",
			"comment varchar(256) not null default '',",
			"repeat varchar(128) not null default ''",
			");")

		_, err = db.Exec(createTable)
		if err != nil {
			log.Fatalf("Cannot create a table with db.Exec(): %s\n", err)
		}
		indexById := "create index id_index on scheduler (id);"
		_, err = db.Exec(indexById)
		if err != nil {
			log.Fatalf("Cannot create an index for table scheduler: %s\n", err)
		}
	}
}
