package tools

import (
	"time"
)

func Time() (string, string) {
	current := time.Now()
	state := "am"

	// check wether the time is past or after morning
	if current.Hour() >= 12 {
		state = "pm"
	}

	date := time.Now().Format("Jan 2, 2006")
	hour := time.Now().Format("03:04" + " " + state)
	return date, hour

}
