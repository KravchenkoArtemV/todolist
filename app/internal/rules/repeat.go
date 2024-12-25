package rules

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var formatTime = "20060102"

// NextDate вычисляет следующую дату в зависимости от правила повторения
func NextDate(nowStr, dateStr, repeat string) (string, error) {
	// Проверки на валидность входных данных
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

	// Обработка правила повторения
	parts := strings.Split(repeat, " ")
	if len(parts) < 1 {
		return "", errors.New("некорректное правило")
	}

	switch parts[0] {
	case "d": // Ежедневное повторение
		if len(parts) != 2 {
			return "", errors.New("некорректное правило")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("некорректное число дней в правиле")
		}
		return calculateNextDate(startDate, now, 0, 0, days), nil

	case "y": // Ежегодное повторение
		return calculateNextDate(startDate, now, 1, 0, 0), nil

	case "w": // Еженедельное повторение
		if len(parts) != 2 {
			return "", errors.New("некорректное правило")
		}
		weekdays, err := parseWeekdays(parts[1])
		if err != nil {
			return "", err
		}
		nextDate, err := calculateNextWeekday(startDate, now, weekdays)
		if err != nil {
			return "", err
		}
		return nextDate, nil

	case "m": // Ежемесячное повторение
		if len(parts) < 2 {
			return "", errors.New("некорректное правило")
		}
		return calculateMonthlyRepeat(parts[1:], startDate, now)

	default:
		return "", errors.New("некорректное правило")
	}
}

// calculateNextDate вычисляет следующую дату с использованием AddDate
func calculateNextDate(startDate, now time.Time, years, months, days int) string {
	nextDate := startDate.AddDate(years, months, days)
	for !nextDate.After(now) {
		nextDate = nextDate.AddDate(years, months, days)
	}
	return nextDate.Format(formatTime)
}

// parseWeekdays парсит список дней недели
func parseWeekdays(weekdays string) ([]int, error) {
	parts := strings.Split(weekdays, ",")
	result := []int{}
	for _, part := range parts {
		day, err := strconv.Atoi(part)
		if err != nil || day < 1 || day > 7 {
			return nil, errors.New("некорректный день недели в правиле")
		}
		result = append(result, day)
	}
	return result, nil
}

// calculateNextWeekday вычисляет ближайший день недели из списка
func calculateNextWeekday(startDate, now time.Time, weekdays []int) (string, error) {
	nextDate := startDate
	for {
		if nextDate.After(now) {
			weekday := int(nextDate.Weekday())
			if weekday == 0 {
				weekday = 7 // Воскресенье — это 7-й день недели
			}
			for _, day := range weekdays {
				if weekday == day {
					return nextDate.Format(formatTime), nil
				}
			}
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}

// calculateMonthlyRepeat вычисляет повторение по дням месяца
func calculateMonthlyRepeat(parts []string, startDate, now time.Time) (string, error) {
	days := []int{}
	months := map[int]bool{}

	// Разбор дней месяца
	for _, part := range strings.Split(parts[0], ",") {
		day, err := strconv.Atoi(part)
		if err != nil || day == 0 || day < -31 || day > 31 {
			return "", errors.New("некорректный день месяца в правиле")
		}
		days = append(days, day)
	}

	// Разбор месяцев (если указаны)
	if len(parts) > 1 {
		for _, part := range strings.Split(parts[1], ",") {
			month, err := strconv.Atoi(part)
			if err != nil || month < 1 || month > 12 {
				return "", errors.New("некорректный месяц в правиле")
			}
			months[month] = true
		}
	}

	// Поиск следующей подходящей даты
	nextDate := startDate
	for {
		if nextDate.After(now) {
			day := nextDate.Day()
			month := int(nextDate.Month())
			for _, d := range days {
				if d == day || (d == -1 && day == lastDayOfMonth(nextDate)) {
					if len(months) == 0 || months[month] {
						return nextDate.Format(formatTime), nil
					}
				}
			}
		}
		nextDate = nextDate.AddDate(0, 0, 1)
	}
}

// lastDayOfMonth возвращает последний день месяца
func lastDayOfMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
}
