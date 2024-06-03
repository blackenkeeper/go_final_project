package config

import "os"

// Setting структура конфигурации, подгружаемой из переменных окружения
type envs struct {
	Port     string
	DBFile   string
	Password string
}

const (
	// DateFormat формат даты, используемый в приложении
	DateFormat = "20060102"
)

var (
	Setting = newSetting()
)

// NewSetting создает новую конфигурацию и загружает значения из переменных окружения
func newSetting() *envs {
	return &envs{
		Port:     getEnv("TODO_PORT", "7540"),
		DBFile:   getEnv("TODO_DBFILE", "scheduler.db"),
		Password: getEnv("TODO_PASSWORD", ""), // Если пароль пуст, оставим поле пустым
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию, если переменная не установлена
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
