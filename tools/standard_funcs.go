package tools

import (
	"log"
	"strconv"
	"time"
)

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalln("⚠ Error while converting to int")
		return 0
	}
	return i
}

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
