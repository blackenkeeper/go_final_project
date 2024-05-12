package logs

import (
	"log"
	"os"
)

func SetupLogging() {
	logFile, err := os.OpenFile("../logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Unable to open log file: %s\n", err)
	}
	log.SetOutput(logFile)
}
