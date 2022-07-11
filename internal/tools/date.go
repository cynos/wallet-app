package tools

import "time"

func DateValidation(start, end string) bool {
	var dateValid bool
	_, a := time.Parse("2006-01-02", start)
	_, b := time.Parse("2006-01-02", end)
	if a == nil && b == nil {
		dateValid = true
	}
	return dateValid
}
