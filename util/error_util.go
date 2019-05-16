package util

import "github.com/golang/glog"

// check and log in error level
func CheckError(tag string, err error) {
	if err != nil {
		glog.Errorf("%s error: %s", tag, err)
	}
}

func CheckFatal(tag string, err error) {
	if err != nil {
		glog.Fatalf("%s error: %s", tag, err)
	}
}
