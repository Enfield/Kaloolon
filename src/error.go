package main

import "strings"

func HandleApiError(err error) bool {
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

func HandleFatalError(err error, message string) {
	if err != nil {
		Error.Fatalf(message + ": %v", err.Error())
	}
}
