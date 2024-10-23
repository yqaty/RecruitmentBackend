package utils

import (
	"log"
	"strconv"
	"time"
)

func TimeParse(ti string) time.Time {
	format := "2006-01-02"
	t, _ := time.Parse(format, ti)
	return t
}

func TimeParseBool(ti string) bool {
	format := "2006-01-02"
	if _, err := time.Parse(format, ti); err != nil {
		log.Println(ti, err)
		return false
	}
	return true
}
func ComPareTimeHour(t1 time.Time, t2 time.Time) bool {
	truncatedTime1 := t1.Truncate(time.Hour)
	truncatedTime2 := t2.Truncate(time.Hour)
	return truncatedTime1.Equal(truncatedTime2)
}

func ConverToLocationTime(t time.Time) (string, error) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	t = t.In(location)
	formattedTime := t.Format("2006年1月2日 Monday 15时04分05秒")
	return formattedTime, nil
}

func GetTimeByString(joinTime string) time.Time {
	year := joinTime[:4]
	season := joinTime[4:]
	month := 6
	if season == "C" {
		month = 9
	} else if season == "A" {
		month = 12
	}
	if y, err := strconv.Atoi(year); err != nil {
		return time.Date(2050, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	} else {
		return time.Date(y, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}
}
