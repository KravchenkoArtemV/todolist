package rules

import (
	"errors"
	"strings"
	"time"
)

var formatTime = "20060102"

// вычисляет следующую дату в зависимости от правила повторения
func NextDate(nowStr, dateStr, repeat string) (string, error) {
	// проверяем, указано ли правило повторения
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// парсим текущую дату из строки
	now, err := time.Parse(formatTime, nowStr)
	if err != nil {
		return "", errors.New("неверный формат текущей даты")
	}

	// парсим начальную дату задачи из строки
	startDate, err := time.Parse(formatTime, dateStr)
	if err != nil {
		return "", errors.New("неверный формат даты задачи")
	}

	// разделяем правило повторения на части
	partsRepeat := strings.Split(repeat, " ")
	if len(partsRepeat) < 1 {
		return "", errors.New("некорректное правило")
	}

	// определяем тип правила повторения по букве
	switch partsRepeat[0] {
	case "d":
		return DayCheck(now, startDate, partsRepeat)
	case "w":
		return WeekCheck(now, startDate, partsRepeat)
	case "y":
		return YearCheck(now, startDate, partsRepeat)
	case "m":
		return MonthCheck(now, startDate, partsRepeat)
	default:
		return "", errors.New("некорректное правило")
	}
}
