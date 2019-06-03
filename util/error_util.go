package util

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// check and log in error level
// This does not stop the program nor does it fail the test case
func CheckError(tag string, err error) {
	if err != nil {
		log.Errorf("%s error: %s", tag, err)
	}
}

func CheckFatal(tag string, err error) {
	if err != nil {
		log.Fatalf("%s error: %s\n", tag, err)
		os.Exit(1)
	}
}
