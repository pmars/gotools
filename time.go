package gotools

import "time"

const (
	TimeFormat = "2006-01-02 15:04:05"
	TimeFormatNoSep = "20060102150405"
	DateFormat = "2006-01-02"
	DateFormatNoSep = "20060102"
)

func GetNowStr() string {
	return time.Now().Format(TimeFormat)
}

func GetNowNoSepStr() string {
	return time.Now().Format(TimeFormatNoSep)
}

func GetDateStr() string {
	return time.Now().Format(DateFormat)
}

func GetDateNoSepStr() string {
	return time.Now().Format(DateFormatNoSep)
}