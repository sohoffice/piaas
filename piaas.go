package piaas

import (
	"flag"
	"github.com/urfave/cli"
	"log"
)

// Below are common option handler.

// Turn on debug logging if -debug was specified
func HandleDebug(c *cli.Context) error {
	if c.Bool("debug") == true {
		log.Println("Debug is on.")
		flag.Lookup("logtostderr").Value.Set("true")
	}
	return nil
}
