package main

import (
	"log"

	"github.com/blackenkeeper/go_final_project/internal/database"
	"github.com/blackenkeeper/go_final_project/internal/logs"
	"github.com/blackenkeeper/go_final_project/internal/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v\n", r)
		}
	}()

	logs.SetupLogging()
	database.SetupDB()
	server.SetupServer()
}
