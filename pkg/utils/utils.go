package utils

import "time"

func DayToDuration(t time.Duration) time.Time {
	year, month, day := time.Now().Date()
	t2 := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return t2.Add(t)
}

func TimePointer(time time.Time) *time.Time {
	return &time
}