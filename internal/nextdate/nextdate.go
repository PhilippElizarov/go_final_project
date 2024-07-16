package nextdate

import (
	"errors"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PhilippElizarov/go_final_project/internal/model"
)

type ByDate []time.Time

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Before(a[j]) }

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("не указано правило повторения")
	}

	date_, err := time.Parse(model.TimeTemplate, date)
	if err != nil {
		return "", err
	}

	s := strings.Split(repeat, " ")

	daysLater := date_

	switch repeat[0] {
	case 'd':
		if len(s) != 2 {
			return "", errors.New("не указан интервал в днях")
		}

		//проверяем корректность введенных дней
		num, err := strconv.Atoi(s[1])
		if err != nil {
			return "", err
		}

		if !(num > 0 && num <= 400) {
			return "", errors.New("неверный диапазон дней")
		}

		//вычисляем новую дату
		for {
			daysLater = daysLater.AddDate(0, 0, num)
			res := daysLater.Compare(now)
			if res == 0 || res == 1 {
				break
			}
		}
	case 'y':
		for {
			daysLater = daysLater.AddDate(1, 0, 0)
			res := daysLater.Compare(now)
			if res == 0 || res == 1 {
				break
			}
		}
	case 'w':
		if len(s) != 2 {
			return "", errors.New("не указаны дни недели")
		}

		weekDays := strings.Split(s[1], ",")

		var weekDaysNums []int

		//проверяем корректность введенных дней недели
		for _, day := range weekDays {
			num, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			if !(num >= 1 && num <= 7) {
				return "", errors.New("неверный диапазон дней недели")
			}
			weekDaysNums = append(weekDaysNums, num)
		}

		var dayOfWeek = map[string]int{
			"Monday":    1,
			"Tuesday":   2,
			"Wednesday": 3,
			"Thursday":  4,
			"Friday":    5,
			"Saturday":  6,
			"Sunday":    7,
		}

		//вычисляем новую дату
		for {
			daysLater = daysLater.AddDate(0, 0, 1)
			res := daysLater.Compare(now)
			weekDayNum := dayOfWeek[daysLater.Weekday().String()]
			if res == 1 && slices.Contains(weekDaysNums, weekDayNum) {
				break
			}
		}
	case 'm':
		if len(s) < 2 || len(s) > 3 {
			return "", errors.New("некорректные параметры повторения")
		}

		days := strings.Split(s[1], ",")

		var daysNums []int
		//проверяем корректность введенных дней
		for _, day := range days {
			num, err := strconv.Atoi(day)
			if err != nil {
				return "", err
			}

			if !(num >= 1 && num <= 31) && num != -1 && num != -2 {
				return "", errors.New("неверный диапазон дней")
			}
			daysNums = append(daysNums, num)
		}
		var months []string
		if len(s) > 2 {
			months = strings.Split(s[2], ",")
		}

		var monthsNums []int
		//проверяем корректность введенных месяцев
		for _, month := range months {
			num, err := strconv.Atoi(month)
			if err != nil {
				return "", err
			}

			if !(num >= 1 && num <= 12) {
				return "", errors.New("неверный диапазон месяцев")
			}
			monthsNums = append(monthsNums, num)
		}

		if len(monthsNums) == 0 {
			monthsNums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		}

		var dateNxt time.Time
		var dates []time.Time
		year, _, _ := daysLater.Date()

		for _, day := range daysNums {
			for _, month := range monthsNums {
				for y := year; y < year+5; y++ {
					if day == -1 || day == -2 {
						firstOfMonth := time.Date(y, time.Month(month), 1, 0, 0, 0, 0, daysLater.Location())
						dateNxt = time.Date(y, time.Month(month), firstOfMonth.AddDate(0, 1, day).Day(), daysLater.Hour(),
							daysLater.Minute(), daysLater.Second(), daysLater.Nanosecond(), daysLater.Location())
					} else {
						dateNxt = time.Date(y, time.Month(month), day, daysLater.Hour(),
							daysLater.Minute(), daysLater.Second(), daysLater.Nanosecond(), daysLater.Location())
						if dateNxt.Month() != time.Month(month) || dateNxt.Day() != day {
							continue
						}
					}
					cmp := dateNxt.Compare(daysLater)
					if cmp == 1 {
						dates = append(dates, dateNxt)
					}
				}
			}
		}

		sort.Sort(ByDate(dates))
		for _, d := range dates {
			cmp := d.Compare(now)
			if cmp == 1 {
				daysLater = d
				break
			}
		}
	default:
		return "", errors.New("неподдерживаемый формат повторения")
	}

	return daysLater.Format(model.TimeTemplate), nil
}
