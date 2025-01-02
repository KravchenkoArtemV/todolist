package rules

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// извлекает дни недели из строки
func executeWeekdays(days string) (map[time.Weekday]bool, error) {
	dayStrings := strings.Split(days, ",") // разделяем строку на части
	weekdaysMap := make(map[time.Weekday]bool)

	for _, dayStr := range dayStrings {
		dayInt, err := strconv.Atoi(dayStr)         // преобразуем строку в число
		if err != nil || dayInt < 1 || dayInt > 7 { // проверяем диапазон от 1 до 7 (пнд-вск)
			return nil, errors.New("некорректный день недели в правиле повторения")
		}
		weekday := time.Weekday((dayInt % 7)) // преобразуем число в день недели остатком от деления
		weekdaysMap[weekday] = true           // добавляем в мапу допустимых значений
	}
	return weekdaysMap, nil
}

// неделя
func WeekCheck(now time.Time, taskDate time.Time, rules []string) (string, error) {
	// проверяем, что правило состоит из двух частей
	if len(rules) != 2 {
		return "", errors.New("неверный формат правила повторения для недели")
	}

	// извлекаем дни недели в виде мапы
	weekdaysMap, err := executeWeekdays(rules[1])
	if err != nil {
		return "", err
	}

	// проходим день за днём, начиная с даты задачи
	for {
		// проверяем входит ли текущий день недели в мапу допустимых значений
		if weekdaysMap[taskDate.Weekday()] {
			if taskDate.After(now) { // если дата задачи позже текущей, возвращаем её
				return taskDate.Format(FormatTime), nil
			}
		}
		// если текущий день не подходит переходим к следующему
		taskDate = taskDate.AddDate(0, 0, 1)
	}
}
