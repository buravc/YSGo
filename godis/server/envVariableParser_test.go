package godis

import (
	"testing"
	"time"
)

func Test_getNumber_NoSpace(t *testing.T) {

	testStr := "12h"
	r := getNumberRegex()

	timeNum, duration, err := getNumber(r, testStr)

	if timeNum != 12 || err != nil || duration != time.Hour {
		t.Fail()
	}

}

func Test_getNumber_WithSpace(t *testing.T) {

	testStr := "12 h"
	r := getNumberRegex()

	timeNum, duration, err := getNumber(r, testStr)

	if timeNum != 12 || err != nil || duration != time.Hour {
		t.Fail()
	}

}
