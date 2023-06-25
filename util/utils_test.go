package util

import (
	"testing"
)

func TestValidateScheduleAdd(t *testing.T) {
	validStrings := []string{"/sa 12:34 5", "/sa 12:34 15"}
	for _, testString := range validStrings {
		if _, err := ValidateSchedule(testString); err != nil {
			t.Errorf(err.Error())
		}
	}
	invalidStrings := []string{"/sa 12:34 xx", "/sa 12:34 ", "/sa 12:3a 1", "/sa  12:34 9", "/sa 25:12 1", "/sa 12:34 0", "/sa 12:34 99"}
	for _, testString := range invalidStrings {
		if _, err := ValidateSchedule(testString); err == nil {
			t.Errorf("ValidateSchedule should have failed for \"%s\"", testString)
		}
	}
}
