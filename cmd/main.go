package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/blackenkeeper/go_final_project/internal/logs"
	"github.com/blackenkeeper/go_final_project/internal/server"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic: %v\n", r)
		}
	}()

	logs.SetupLogging()
	server.SetupServer()
}
