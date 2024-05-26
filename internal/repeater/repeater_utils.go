package repeater

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Функция для проверки года на високосность
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// Функция для нахождения следующего месяца, в котором должна повторяться задача
func nearestMonthFinder(date time.Time, monthsLegit map[int]bool) time.Time {
	log.Println("Запуск функции nearestMonthFinder для нахождения следующего месяца повторения задачи")
	nextMonth := int(date.Month()) + 1
	if nextMonth > 12 {
		log.Println("Текущий месяц - декабрь, поэтому следующим для проверки валидности выбираем январь")
		nextMonth = 1
		date = date.AddDate(1, 0, 0)
	}

	var found bool
	log.Printf("Поиск следующего месяца после %d для повторения задачи\n", nextMonth)
	for i := nextMonth; i < len(monthsLegit)+1; i++ {
		if monthsLegit[i] {
			log.Println("Значение найдено! Следующий валидный для повторения месяц:", time.Month(i))
			found = true
			date = time.Date(date.Year(), time.Month(i), 1, 0, 0, 0, 0, time.UTC)
			break
		}
	}

	if !found {
		log.Println("Валидных месяцев в этом году не осталось, ищем первый валидный месяц в году")
		for i := 1; i < nextMonth; i++ {
			if monthsLegit[i] {
				log.Println("Найден! Выбираем его для повторения, увеличив счётчик лет на единицу")
				date = time.Date(date.Year()+1, time.Month(i), 1, 0, 0, 0, 0, time.UTC)
				break
			}
		}
	}

	log.Println("Проверяем год на високосность")
	if !isLeapYear(date.Year()) {
		log.Println("Год не високосный")
		monthsDays[2] = 28
	} else {
		log.Println("Год високосный")
		monthsDays[2] = 29
	}

	log.Println("Успешное окончание работы функции nearestMonthFinder")
	return date
}

// Костыль, на котором держится логика для обработки повторения задач по месяцам.
// Функция перебирает все дни для повторения в месяце из repeatDays и добавляет числовые
// значения в массив чисел repeatDaysInt. После записи всех значений массив сортируется по возрастанию.
// Если текущий день в задаче есть в repeatDays, то записывает в перемнную isTodayInRepeatList
// значение true, иначе false. Сегодняшний день добавляется в массив repeatDaysInt в любом случае.
//
// Функция возвращает числовой массив дней повторения в месяце, ошибку или nil, если её нет,
// и значение переменной isTodayInRepeatList
func allRepeatDays(taskDate time.Time, repeatDays []string) ([]int, error, bool) {
	log.Println("Начало работы функции allRepeatDays")
	repeatDaysInt := []int{taskDate.Day()}
	isTodayInRepeatList := false

	log.Println("Попытка парсинга значений массива строк repeatDays в числовые значения")
	for _, day := range repeatDays {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt == 0 || dayInt < -2 {
			log.Println("Переданный параметр не является числом, равен нулю или меньше, чем -2")
			return nil,
				errors.New("переданный параметр не является числом, равен нулю или меньше, чем -2"), false
		}
		if dayInt > 31 {
			log.Println("Введенное значение дней за пределами диапазона 1-31")
			return nil,
				errors.New("введенное значение дней за пределами диапазона 1-31"), false
		}

		if dayInt < 0 {
			log.Printf("Значение %d отсчитывает дни с конца месяца и преобразуется в ", dayInt)
			dayInt += monthsDays[int(taskDate.Month())] + 1
			log.Print(dayInt, "\n")
		}

		log.Println("Проверяем, является ли сегодняшний по taskDate день пригодным для повторения задачи.")
		log.Println("Если не выведется \"Является\" в следующей строке, значит, не является")
		if dayInt == repeatDaysInt[0] {
			log.Println("Является")
			isTodayInRepeatList = true
			continue
		}

		log.Printf("Добавляем значение %d в массив repeatDaysInt\n", dayInt)
		repeatDaysInt = append(repeatDaysInt, dayInt)
	}

	log.Println("Сортируем массив repeatDaysInt после добавления и преобразования всех значений из ",
		"repeatDays")
	sort.Slice(repeatDaysInt, func(i, j int) bool {
		return repeatDaysInt[i] < repeatDaysInt[j]
	})

	log.Println("Успешное окончание работы функции allRepeatDays")
	return repeatDaysInt, nil, isTodayInRepeatList
}

// monthsThirdParamChecker проверяет третий параметр правила повторения задачи, если используется
// повторение по месяцам.
func monthsThirdParamChecker(mLegit *map[int]bool, repeatRule []string) error {
	log.Println("Начало работы функции monthsThirdParamChecker")
	monthsLegit := *mLegit
	if len(repeatRule) == 3 {
		monthsNumber := strings.Split(repeatRule[2], ",")
		log.Println("В правиле повторения заданы месяцы, в которые должна повторяться задача:",
			monthsNumber)
		log.Println("Проверяем эти значения месяцев на валидность и, в случае успеха, отмечаем только их",
			"пригодным для повторения")
		for _, month := range monthsNumber {
			monthInt, err := strconv.Atoi(month)

			if err != nil || monthInt < 1 || monthInt > 12 {
				log.Printf("Значение %s за пределами допустимого числового диапазона (1-12)\n", month)
				return errors.New("введённое число за пределами допустимого значения месяца (1-12)")
			}
			if _, exists := monthsLegit[monthInt]; exists {
				monthsLegit[monthInt] = true
			}
		}
	} else {
		log.Println("В правиле повторения не указаны конкретные месяцы, значит, сделаем все месяцы валидными")
		for i := 1; i < len(monthsLegit)+1; i++ {
			monthsLegit[i] = true
		}
	}

	log.Println("Успешное окончание работы функции monthsThirdParamChecker")
	return nil
}
