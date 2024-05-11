package main

import (
	"log"
	//"github.com/blackenkeeper/go_final_project/logs"
	//"github.com/blackenkeeper/go_final_project/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v\n", r)
		}
	}()

	logs.setupLogging()
	server.setupDB()
	server.setupServer()
}
