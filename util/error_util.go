package util

import (
	"fmt"
	"github.com/golang/glog"
	"os"
)

// check and log in error level
func CheckError(tag string, err error) {
	if err != nil {
		glog.Errorf("%s error: %s", tag, err)
	}
}

func CheckFatal(tag string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s error: %s\n", tag, err)
		os.Exit(1)
	}
}
