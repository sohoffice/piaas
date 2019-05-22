// +build darwin

package piaas

import (
	log "github.com/sirupsen/logrus"
)

func NewMonitor(startDir string) Monitor {
	log.Infof("Use fsevent monitor.")
	fsm := NewFSEventMonitor(startDir)
	return &fsm
}
