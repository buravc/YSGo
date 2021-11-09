package godis

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

func getNumberRegex() *regexp.Regexp {
	return regexp.MustCompile(`^([0-9]+?)\s*(?:(s|h|m))?$`)
}

func parseString(str string) time.Duration {
	r := getNumberRegex()
	timeNum, duration, err := getNumber(r, str)

	if err != nil {
		log.Fatal("Env variable could not be parsed!")
	}

	return time.Duration(timeNum) * duration
}

func getNumber(r *regexp.Regexp, str string) (int16, time.Duration, error) {
	found := r.FindStringSubmatch(str)
	len := len(found)

	if len != 3 {
		return 0, 0, &strconv.NumError{}
	}

	foundStr := found[1]
	res, err := strconv.ParseInt(foundStr, 0, 16)

	var duration time.Duration

	switch found[2] {
	case "h":
		duration = time.Hour
	case "m":
		duration = time.Minute
	case "s":
		duration = time.Second
	default:
		return 0, 0, &strconv.NumError{}
	}

	return int16(res), duration, err
}
