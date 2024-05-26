package utils

import (
	"strings"
)

// Функция, меняющая путь текущей директории cmd на корневую (возвращает на уровень выше)
func CmdPathChecker(filepath string) string {
	if strings.HasSuffix(filepath, "cmd") {
		filepath = strings.Trim(filepath, "cmd")
	} else {
		filepath += "/"
	}

	return filepath
}
