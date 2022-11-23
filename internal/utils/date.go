package utils

import "time"

func GetToday() string {
	return time.Now().Format("2006-01-02")
}
