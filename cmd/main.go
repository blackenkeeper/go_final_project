package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/blackenkeeper/go_final_project/internal/logs"
	"github.com/blackenkeeper/go_final_project/internal/server"
)

func main() {
	logs.SetupLogging()

	err := server.SetupServer()
	if err != nil {
		log.Error(err)
		return
	}
}
