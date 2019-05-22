// +build !darwin

package piaas

import (
	log "github.com/sirupsen/logrus"
)

func NewMonitor(startDir string) Monitor {
	log.Infof("Use fsnotify recursive monitor.")
	rm := NewRecursiveMonitor(startDir)
	return &rm
}
