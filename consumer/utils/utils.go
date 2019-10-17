package utils

import "log"

// FailOnError - utility function for logging errors and exiting the process
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// LogOnError - utility function for logging errors without exiting the process
func LogOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}
