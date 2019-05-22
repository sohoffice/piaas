package piaas

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Below are common option handler.

// Turn on debug logging if -debug was specified
func HandleDebug(c *cli.Context) error {
	if c.Bool("debug") == true {
		log.Println("Debug is on.")
		log.SetLevel(log.DebugLevel)
	}
	return nil
}
