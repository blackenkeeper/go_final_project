package main

import (
	"log"
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
