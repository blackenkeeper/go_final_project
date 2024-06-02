package repeater

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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

func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		log.WithError(err).Error("Неверный формат даты")
		return "", err
	}

	repeatRule := strings.Fields(repeat)
	if len(repeatRule) < 1 {
		err := errors.New("неверный формат правила повторения задачи")
		log.Error(err)
		return "", err
	}

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
		err := errors.New("не соответствует ни одному из правил повторения")
		log.Error(err)
		return "", err
	}
}

func yearRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) != 1 {
		err := errors.New("для параметра правила 'y' нельзя указать число, только сам параметр")
		log.Error(err)
		return "", err
	}

	if taskDate.After(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}

	for !taskDate.After(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}

	return taskDate.Format("20060102"), nil
}

func dayRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) != 2 {
		err := errors.New("неверный формат правила для 'd': укажите ОДНО число в интервале от 1 до 400")
		log.Error(err)
		return "", err
	}

	days, err := strconv.Atoi(repeatRule[1])
	if err != nil {
		log.WithError(err).Error("Ошибка парсинга значения в число")
		return "", err
	}

	if days < 1 || days > 400 {
		err := errors.New("значение для 'd' за пределами допустимого диапазона от 1 до 400")
		log.Error(err)
		return "", err
	}

	if taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}

	for !taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}

	return taskDate.Format("20060102"), nil
}

func weekRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) != 2 {
		err := errors.New("не указан номер дня недели или неверный формат правила")
		log.Error(err)
		return "", err
	}

	if taskDate.Before(now) {
		taskDate = now
	}

	days := strings.Split(repeatRule[1], ",")

	for _, day := range days {
		if _, exists := weekdays[day]; !exists {
			err := errors.New("введённое значение не является числом или за пределами диапазона 1-7")
			log.Error(err)
			return "", err
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
		err := errors.New("не указаны дни месяца для повторения задачи или превышено допустимое количество параметров")
		log.Error(err)
		return "", err
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
		log.WithError(err).Error("Ошибка обработки дней повторения")
		return "", err
	}

	err = monthsThirdParamChecker(&monthsLegit, repeatRule)
	if err != nil {
		log.WithError(err).Error("Ошибка проверки третьего параметра")
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
			if !monthsLegit[currentMonth] || closestDay == repeatDaysInt[0] || closestDay > monthsDays[currentMonth] {
				taskDate = nearestMonthFinder(taskDate, monthsLegit)
				break
			}
			found = true
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), closestDay, 0, 0, 0, 0, time.UTC)
			break
		}
	}

	if !found {
		repeatDaysInt, err, firstDayRepeatable := allRepeatDays(taskDate, repeatDays)
		if err != nil {
			return "", err
		}

		if firstDayRepeatable {
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), repeatDaysInt[0], 0, 0, 0, 0, time.UTC)
		} else {
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), repeatDaysInt[1], 0, 0, 0, 0, time.UTC)
		}
	}

	return taskDate.Format("20060102"), nil
}
