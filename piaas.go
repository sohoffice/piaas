package piaas

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

// Prepare common flags for commands.
func PrepareCommonFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "Specify the piaas config file name and path. Default to piaasconfig.yml",
			Value: "piaasconfig.yml",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Print debug messages",
		},
	}
}

// Below are common option handler.

// Turn on debug logging if -debug was specified
func HandleDebug(c *cli.Context) error {
	if c.Bool("debug") == true {
		log.Println("Debug is on.")
		log.SetLevel(log.DebugLevel)
	}
	return nil
}

func HandleExitSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, os.Kill)

	go func() {
		sig := <-signalCh
		// no longer accept other signals.
		close(signalCh)
		log.Debugf("Signal received: %s", sig)

		// notify all observers.
		for _, obs := range exitObservers {
			obs(sig)
		}

	}()
}

var exitObservers []func(os.Signal)

// Add observer to listen to exit signal.
// Set preferredFront to true to make sure the observer is ranked in front.
func SubscribeExitSignal(obs func(os.Signal), preferredFront bool) {
	if preferredFront {
		log.Debugf("Add an exit observer in the first position")
		exitObservers = append([]func(os.Signal){obs}, exitObservers...)
	} else {
		log.Debugf("Add an exit observer in the last position")
		exitObservers = append(exitObservers, obs)
	}
}
