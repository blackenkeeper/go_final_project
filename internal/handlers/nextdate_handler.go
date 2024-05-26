package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/repeater"
)

// Обработчик пути /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запуск обработчика NextDateHandler для пути /api/nextdate")

	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	log.Printf("Получениe параметров из запроса: now - %s, date - %s, repeat - %s",
		nowParam, dateParam, repeatParam)

	// Проверяем наличие всех параметров
	if nowParam == "" || dateParam == "" || repeatParam == "" {
		log.Println("Не указаны значения для всех параметров")
		http.Error(w, "Не указаны все параметры", http.StatusBadRequest)
		return
	}

	log.Println("Парсинг nowParam в объект даты")
	nowTime, err := time.Parse("20060102", nowParam)
	if err != nil {
		log.Println("Ошибка парсинга параметра now: ", err.Error())
		http.Error(w, "Ошибка парсинга параметра now: "+err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := repeater.NextDate(nowTime, dateParam, repeatParam)
	if err != nil {
		http.Error(w, "Ошибка при получении следующей даты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Обработчик NextDateHandler отработал успешно")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, nextDate)
}
