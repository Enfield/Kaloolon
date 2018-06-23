package main

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		Error.Fatalf(message + ": %v", err.Error())
	}
}