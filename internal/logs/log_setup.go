package logs

import (
	log "github.com/sirupsen/logrus"
)

// Функция для настройки работы логов
func SetupLogging() {
	// Установка формата вывода
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	// Установите уровень логирования
	log.SetLevel(log.DebugLevel)
}
