package repeater

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

var weekdays = map[string]int{
	"1": int(time.Monday),
	"2": int(time.Tuesday),
	"3": int(time.Wednesday),
	"4": int(time.Thursday),
	"5": int(time.Friday),
	"6": int(time.Saturday),
	"7": int(time.Sunday),
}

// Словарь "месяц: количество дней в нём"
var monthsDays = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

// NextDate возвращает следующую дату по заданному правилу в repeat. Формат правил:
//
// "d 1" - повторять задачу через день (то есть каждый день). Значение можно указать от 1 до 400.
//
// "y" - повторять задачу каждый год. Указывать нужно только этот символ, перенести можно ровно
// на год от текущей даты задачи.
//
// "w 3,4,5" - повторять в 3, 4 и 5 дни недели (среда, четверг, пятница). Дни указываются через запятую,
// без проблеов.
//
// "m 2,3,-2 3,6,9,12" - повторять задачу 2 и 3 числа, а также в предпоследний день месяца
// в марте, июне, сентябре и декбаре. Также можно указывать -1 для последнего дня месяца.
// Третий параметр (параметры разделены пробелами) с конкретными месяцами опциональный, и если не указан,
// то задача будет повторяться во все месяцы, в которых есть нужный день повторения
// (например, если указать "m 30", то в феврале точно не будет повторения).
func NextDate(now time.Time, date string, repeat string) (string, error) {
	log.Println("Начало работы функции NextDate")
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		log.Println("Неверный формат даты:", err)
		return "", err
	}

	log.Printf("Разбивка правила \"%s\" на составляющие\n", repeat)
	repeatRule := strings.Fields(repeat)
	if len(repeatRule) < 1 {
		log.Println("Неверный формат правила повторения задачи")
		return "", errors.New("неверный формат правила повторения задачи")
	}

	log.Println("Передача в нужный обработчик правила для значения", repeatRule[0])
	switch repeatRule[0] {
	case "y":
		return yearRule(now, taskDate, repeatRule)
	case "d":
		return dayRule(now, taskDate, repeatRule)
	case "w":
		return weekRule(now, taskDate, repeatRule)
	case "m":
		return monthRule(now, taskDate, repeatRule)
	default:
		return "", errors.New("не соответствует ни одному из правил повторения")
	}
}

// Обработчик правила повторения по годам ('y')
func yearRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	log.Println("Запуск функции yearRule (правило 'y') для повторения по годам")
	if len(repeatRule) != 1 {
		log.Println("Для параметра правила 'y' нельзя указать число, только сам параметр")
		return "", errors.New("для параметра правила 'y' нельзя указать число, только сам параметр")
	}

	if taskDate.After(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}

	for !taskDate.After(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}

	log.Println("Функция yearRule успешо завершила свою работу")
	return taskDate.Format("20060102"), nil
}

// Обработчик правила повторения по дням ('d')
func dayRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	log.Println("Запуск функции dayRule (правило 'd') для повторения по дням")
	if len(repeatRule) < 2 {
		log.Println("Не указано количество дней для повторения задачи")
		return "", errors.New("не указано количество дней для повторения задачи")
	}

	if len(repeatRule) > 2 {
		log.Println("Невереный формат правила для 'd': укажите ОДНО число в интервале от 1 до 400")
		return "",
			errors.New("невереный формат правила для 'd': укажите ОДНО число в интервале от 1 до 400")
	}

	log.Printf("Парсинг значения %s в число\n", repeatRule[1])
	days, err := strconv.Atoi(repeatRule[1])
	if err != nil {
		log.Println(err)
		return "", err
	}

	if days < 1 || days > 400 {
		log.Println("Значение для 'd' за пределами допустимого диапазона от 1 до 400")
		return "", errors.New("значение для 'd' за пределами допустимого диапазона от 1 до 400")
	}

	if taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}

	for !taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}

	log.Println("Функция dayRule успешно завершила работу")
	return taskDate.Format("20060102"), nil
}

// Обработчик правила повторения по неделям ('w') (логи допишу позже, обещаю)
func weekRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) < 2 || len(repeatRule) > 2 {
		return "", errors.New("не указан номер дня недели или неверный формат правила")
	}

	if taskDate.Before(now) {
		taskDate = now
	}

	days := strings.Split(repeatRule[1], ",")

	for _, day := range days {
		if _, exists := weekdays[day]; !exists {
			return "", errors.New("введённое значение не является числом или за пределами диапазона 1-7")
		}
	}

	sort.Slice(days, func(i, j int) bool {
		return days[i] < days[j]
	})

	findAtFirstIter := true
	for i := 0; i < len(days); i++ {
		todayWeekday := taskDate.Weekday()
		nextWeekday := time.Weekday(weekdays[days[i]])

		if todayWeekday == nextWeekday {
			if findAtFirstIter {
				if len(days) == 1 {
					taskDate = taskDate.AddDate(0, 0, 7)
					break
				}
				neededWeekday := weekdays[days[0]]
				if i+1 != len(days) {
					neededWeekday = weekdays[days[i+1]]
				}
				for taskDate.Weekday() != time.Weekday(neededWeekday) {
					taskDate = taskDate.AddDate(0, 0, 1)
				}
			}
			break
		}

		if i+1 == len(days) {
			findAtFirstIter = false
			taskDate = taskDate.AddDate(0, 0, 1)
			i = -1
		}
	}

	return taskDate.Format("20060102"), nil
}

// Обработчик правила повторения по месяцам ('m') (логи допишу позже, обещаю)
func monthRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	var monthsLegit = map[int]bool{
		1:  false,
		2:  false,
		3:  false,
		4:  false,
		5:  false,
		6:  false,
		7:  false,
		8:  false,
		9:  false,
		10: false,
		11: false,
		12: false,
	}

	if len(repeatRule) < 2 || len(repeatRule) > 3 {
		return "",
			errors.New("не указаны дни месяца для повторения задачи или превышено " +
				"допустимое количество параметров")
	}
	if taskDate.Before(now) {
		taskDate = now
	}
	if isLeapYear(taskDate.Year()) {
		monthsDays[2] = 29
	}

	repeatDays := strings.Split(repeatRule[1], ",")
	repeatDaysInt, err, _ := allRepeatDays(taskDate, repeatDays)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = monthsThirdParamChecker(&monthsLegit, repeatRule)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var found bool
	for i := 0; i < len(repeatDaysInt); i++ {
		if repeatDaysInt[i] == taskDate.Day() {
			closestDay := repeatDaysInt[0]
			currentMonth := int(taskDate.Month())

			if i+1 != len(repeatDaysInt) {
				closestDay = repeatDaysInt[i+1]
			}
			if !monthsLegit[currentMonth] || closestDay == repeatDaysInt[0] ||
				closestDay > monthsDays[currentMonth] {
				taskDate = nearestMonthFinder(taskDate, monthsLegit)
				break
			}
			found = true
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), closestDay,
				0, 0, 0, 0, time.UTC)
			break
		}
	}

	if !found {
		repeatDaysInt, err, firstDayRepeatable := allRepeatDays(taskDate, repeatDays)
		if err != nil {
			return "", err
		}

		if firstDayRepeatable {
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), repeatDaysInt[0],
				0, 0, 0, 0, time.UTC)
		} else {
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), repeatDaysInt[1],
				0, 0, 0, 0, time.UTC)
		}

	}

	return taskDate.Format("20060102"), nil
}
