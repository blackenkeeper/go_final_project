package repeater

import (
	"errors"
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
	nextMonth := int(date.Month()) + 1
	if nextMonth > 12 {
		nextMonth = 1
		date = date.AddDate(1, 0, 0)
	}

	var found bool
	for i := nextMonth; i < len(monthsLegit)+1; i++ {
		if monthsLegit[i] {
			found = true
			date = time.Date(date.Year(), time.Month(i), 1, 0, 0, 0, 0, time.UTC)
			break
		}
	}

	if !found {
		for i := 1; i < nextMonth; i++ {
			if monthsLegit[i] {
				date = time.Date(date.Year()+1, time.Month(i), 1, 0, 0, 0, 0, time.UTC)
				break
			}
		}
	}

	if !isLeapYear(date.Year()) {
		monthsDays[2] = 28
	} else {
		monthsDays[2] = 29
	}

	return date
}

// Костыль, на котором держится логика для обработки повторения задач по месяцам.
func allRepeatDays(taskDate time.Time, repeatDays []string) ([]int, error, bool) {
	repeatDaysInt := []int{taskDate.Day()}
	isTodayInRepeatList := false

	for _, day := range repeatDays {
		dayInt, err := strconv.Atoi(day)
		if err != nil || dayInt == 0 || dayInt < -2 {
			return nil, errors.New("переданный параметр не является числом, равен нулю или меньше, чем -2"), false
		}
		if dayInt > 31 {
			return nil, errors.New("введенное значение дней за пределами диапазона 1-31"), false
		}

		if dayInt < 0 {
			dayInt += monthsDays[int(taskDate.Month())] + 1
		}

		if dayInt == repeatDaysInt[0] {
			isTodayInRepeatList = true
			continue
		}

		repeatDaysInt = append(repeatDaysInt, dayInt)
	}

	sort.Slice(repeatDaysInt, func(i, j int) bool {
		return repeatDaysInt[i] < repeatDaysInt[j]
	})

	return repeatDaysInt, nil, isTodayInRepeatList
}

// monthsThirdParamChecker проверяет третий параметр правила повторения задачи, если используется
// повторение по месяцам.
func monthsThirdParamChecker(mLegit *map[int]bool, repeatRule []string) error {
	monthsLegit := *mLegit
	if len(repeatRule) == 3 {
		monthsNumber := strings.Split(repeatRule[2], ",")
		for _, month := range monthsNumber {
			monthInt, err := strconv.Atoi(month)
			if err != nil || monthInt < 1 || monthInt > 12 {
				return errors.New("введённое число за пределами допустимого значения месяца (1-12)")
			}
			if _, exists := monthsLegit[monthInt]; exists {
				monthsLegit[monthInt] = true
			}
		}
	} else {
		for i := 1; i < len(monthsLegit)+1; i++ {
			monthsLegit[i] = true
		}
	}

	return nil
}
