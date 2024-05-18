package utils

import (
	"strings"
)

func CmdPathChecker(filepath string) string {
	if strings.HasSuffix(filepath, "cmd") {
		filepath = strings.Trim(filepath, "cmd")
	} else {
		filepath += "/"
	}

	return filepath
}
