package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blackenkeeper/go_final_project/internal/repeater"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatParam := r.URL.Query().Get("repeat")

	// Проверяем наличие всех параметров
	if nowParam == "" || dateParam == "" || repeatParam == "" {
		http.Error(w, "Не указаны все параметры", http.StatusBadRequest)
		return
	}

	// Парсим nowParam в формат даты
	nowTime, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "Ошибка парсинга параметра now: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Вызываем функцию NextDate
	nextDate, err := repeater.NextDate(nowTime, dateParam, repeatParam)
	if err != nil {
		http.Error(w, "Ошибка при получении следующей даты: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, nextDate)
}
