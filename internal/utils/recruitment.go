package utils

import (
	"errors"
	"regexp"
)

var name = map[string]string{
	"S": "春季招新",
	"C": "夏令营招新",
	"A": "秋季招新",
	"春": "春季招新",
	"夏": "夏令营招新",
	"秋": "秋季招新",
}

func ConvertRecruitmentName(title string) string {
	year := title[:4]
	suffix := name[title[4:]]
	if suffix == "" {
		return title
	}
	return year + suffix
}

func CheckNameValid(name string) error {
	if len(name) != 5 {
		return errors.New("recruitment name is invalid, correct format like \"2023S/A/C\" (S是春季招新，C是夏令营，A是秋季招新)")
	}
	_, err := regexp.MatchString(`^\d{4}[SAC]$`, name)
	if err != nil {
		return errors.New(err.Error() + "recruitment name is invalid, correct format like \"2023S/A/C\" (S是春季招新，C是夏令营，A是秋季招新)")
	}
	return nil
}
