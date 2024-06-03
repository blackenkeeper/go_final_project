package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/config"
	"github.com/blackenkeeper/go_final_project/internal/repeater"
	log "github.com/sirupsen/logrus"
)

// Обработчик пути /api/nextdate
func (h *Handler) NextDateHandler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Запуск обработчика NextDateHandler для пути /api/nextdate")

	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	// Проверяем наличие всех параметров
	if nowParam == "" || dateParam == "" || repeatParam == "" {
		log.Warn("Не указаны значения для всех параметров")
		http.Error(w, "Не указаны все параметры", http.StatusBadRequest)
		return
	}

	nowTime, err := time.Parse(config.DateFormat, nowParam)
	if err != nil {
		log.Warn("Ошибка парсинга параметра now:", err)
		http.Error(w, "Ошибка парсинга параметра now: "+err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := repeater.NextDate(nowTime, dateParam, repeatParam)
	if err != nil {
		log.Error("Ошибка при получении следующей даты:", err)
		http.Error(w, "Ошибка при получении следующей даты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Обработчик NextDateHandler отработал успешно")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, nextDate)
}
