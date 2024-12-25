package rules

import (
	"time"
)

// год
func YearCheck(now, startDate time.Time, parts []string) (string, error) {
	// +год к стартовой дате
	nextDate := startDate.AddDate(1, 0, 0)
	// пока текущая дата больше плюсуем год
	for nextDate.Before(now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format(formatTime), nil
}
