package logs

import (
	"log"
	"os"

	"github.com/blackenkeeper/go_final_project/internal/utils"
)

// Функция для настройки работы логов
func SetupLogging() {
	filepath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	filepath = utils.CmdPathChecker(filepath)
	logFile, err := os.OpenFile(filepath+"logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Unable to open log file: %s\n", err)
	}
	log.SetOutput(logFile)
}
