package utils

import (
	"fmt"
	"net/http"
	"os"
)

//Check400M allows to check for error 400
func Check400M(w http.ResponseWriter) func(string, string) bool {
	return func(x string, msg string) bool {
		if x == "" {
			http.Error(w, fmt.Sprintf("[400]- Bad request: %s is missing", msg), http.StatusBadRequest)
			return true
		}
		return false

	}
}

//Check500 allows to check for error 500
func Check500(w http.ResponseWriter) func(error) bool {
	return func(err error) bool {
		if err != nil {
			http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
			return true
		}
		return false
	}
}

//Exit struct of  error code
type Exit struct{ Code int }

// exit code handler
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}
