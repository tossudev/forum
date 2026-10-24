package timeformat

import (
	"time"
)

const (
	TimeFormat         string = "20060102T150405"
	TimeFormatFrontEnd string = "January 02 2006 15:04"
)

func FormatTime(timeString string, loc *time.Location) (string, error) {
	time, err := time.Parse(TimeFormat, timeString)
	if err != nil {
		return "", err
	}
	// location = users location
	time = time.In(loc)
	formattedTime := time.Format(TimeFormatFrontEnd)
	return formattedTime, nil
}
