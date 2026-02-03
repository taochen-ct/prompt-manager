package common

import "time"

const DateTimeLayout = "2006-01-02 15:04:05"

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return t.String()
	}
	return t.Format(DateTimeLayout)
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(DateTimeLayout, s)
}
