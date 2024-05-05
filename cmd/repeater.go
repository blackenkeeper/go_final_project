package main

import (
	"errors"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

var weekDays = map[string]int{
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
		log.Println("Неверный формат даты:", err)
		return "", err
	}

	repeatRule := strings.Fields(repeat)
	if len(repeatRule) < 1 {
		return "", errors.New("неверный формат правила повторения задачи")
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
		return "", errors.New("не соответствует ни одному из правил повторения")
	}
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func yearRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) != 1 {
		return "", errors.New("для параметра правила 'y' нельзя указать число, только сам параметр")
	}

	if taskDate.Equal(now) || taskDate.After(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}

	for taskDate.Before(now) {
		taskDate = taskDate.AddDate(1, 0, 0)
	}

	return taskDate.Format("20060102"), nil
}

func dayRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) < 2 {
		return "", errors.New("не указано количество дней для повторения задачи")
	}

	if len(repeatRule) > 2 {
		return "",
			errors.New("невереный формат правила для 'd': укажите ОДНО число в интервале от 1 до 400")
	}

	days, err := strconv.Atoi(repeatRule[1])
	if err != nil {
		return "", err
	}

	if days < 1 || days > 400 {
		return "", errors.New("значение для 'd' за пределами допустимого диапазона от 1 до 400")
	}

	if taskDate.Equal(now) || taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}

	for taskDate.Before(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	}

	return taskDate.Format("20060102"), nil
}

func weekRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) < 2 || len(repeatRule) > 2 {
		return "", errors.New("не указан номер дня недели или неверный формат правила")
	}

	days := strings.Split(repeatRule[1], ",")
	for _, day := range days {
		if _, exists := weekDays[day]; !exists {
			return "", errors.New("введённое значение не является числом или за пределами диапазона 1-7")
		}
	}

	for i := 0; ; {
		var found bool
		if time.Weekday(weekDays[days[i]]) == taskDate.Weekday() {
			found = true
			if taskDate.After(now) {
				break
			}
			taskDate = taskDate.AddDate(0, 0, 7)
		}

		if !found {
			i++
			if i == len(days) {
				taskDate = taskDate.AddDate(0, 0, 1)
				i = 0
			}
		}
	}

	return taskDate.Format("20060102"), nil
}

func monthRule(now, taskDate time.Time, repeatRule []string) (string, error) {
	if len(repeatRule) < 2 || len(repeatRule) > 3 {
		return "",
			errors.New("не указаны дни месяца для повторения задачи или превышено " +
				"допустимое количество параметров")
	}
	if taskDate.Before(now) {
		taskDate = now
	}
	repeatDays := strings.Split(repeatRule[1], ",")

	repeatDaysInt := []int{taskDate.Day()}
	for _, day := range repeatDays {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt < -2 {
			return "", errors.New("переданный параметр не является числом или меньше, чем -2")
		}

		if dayInt < 0 {
			dayInt += monthsDays[int(taskDate.Month())]
		}

		if dayInt == repeatDaysInt[0] {
			continue
		}

		repeatDaysInt = append(repeatDaysInt, dayInt)
	}

	sort.Slice(repeatDaysInt, func(i, j int) bool {
		return repeatDaysInt[i] < repeatDaysInt[j]
	})

	for i := 0; i < len(repeatDaysInt); i++ {
		if repeatDaysInt[i] == taskDate.Day() {
			closestDay := repeatDaysInt[0]
			if i+1 < len(repeatDaysInt) {
				closestDay = repeatDaysInt[i+1]
			}
			taskDate = time.Date(taskDate.Year(), taskDate.Month(), closestDay,
				0, 0, 0, 0, time.Local)
		}
	}

	return taskDate.Format("20060102"), nil
}
