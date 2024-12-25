package rules

import (
	"errors"
	"strings"
	"time"
)

var formatTime = "20060102"

// вычисляет следующую дату в зависимости от правила повторения
func NextDate(nowStr, dateStr, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	now, err := time.Parse(formatTime, nowStr)
	if err != nil {
		return "", errors.New("неверный формат текущей даты")
	}

	startDate, err := time.Parse(formatTime, dateStr)
	if err != nil {
		return "", errors.New("неверный формат даты задачи")
	}

	parts := strings.Split(repeat, " ")
	if len(parts) < 1 {
		return "", errors.New("некорректное правило")
	}

	switch parts[0] {
	case "d":
		return DayCheck(now, startDate, parts)
	case "w":
		return WeekCheck(now, startDate, parts)
	case "y":
		return YearCheck(now, startDate, parts)
	case "m":
		return MonthCheck(now, startDate, parts)
	default:
		return "", errors.New("некорректное правило")
	}
}
