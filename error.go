package main

import "strings"

func handleApiError(err error) bool {
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "quotaExceeded"){
			Error.Fatalf(msg+": %v", err.Error())
		}
		Error.Printf(msg+": %v", err.Error())
		return false
	}
	return true
}

func handleError(err error, message string) {
	if err != nil {
		Error.Fatalf(message + ": %v", err.Error())
	}
}
