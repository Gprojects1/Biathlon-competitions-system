package utils

import (
	"fmt"
	"time"
)

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("15:04:05.000", timeStr)
}

func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", h, m, s, ms)
}
