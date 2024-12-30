package rules

import (
	"errors"
	"strconv"
	"time"
)

// день
func DayCheck(now, startDate time.Time, parts []string) (string, error) {
	// проверка правила на состав только из 2 частей и преобразование в число второй части
	if len(parts) != 2 {
		return "", errors.New("некорректное правило для ежедневного повторения")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 || days > 400 {
		return "", errors.New("некорректное число дней в правиле")
	}
	nextDate := startDate.AddDate(0, 0, days)
	// пока текущая дата больше плюсуем дни
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(0, 0, days)
	}
	return nextDate.Format(FormatTime), nil
}
