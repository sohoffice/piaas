package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas"
	"github.com/sohoffice/piaas/app"
	"github.com/sohoffice/piaas/sync"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"time"
)

var version string

func main() {
	// Override flag.CommandLine to prevent extraneous flag from stopping the execution. It should be managed by urfave/cli.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	// Do not print any output from flag module.
	flag.CommandLine.SetOutput(ioutil.Discard)
	// We still would like the service of flag
	flag.Parse()
	// default log level is INFO
	log.SetLevel(log.InfoLevel)

	cliApp := cli.NewApp()
	cliApp.Name = "Piaas"
	cliApp.Description = "Increase productivity of developer by leveraging computing power of multiple machines."
	cliApp.HelpName = "piaas"
	cliApp.Authors = []cli.Author{
		{
			Name:  "Douglas Liu",
			Email: "douglas@sohoffice.com",
		},
	}
	cliApp.Version = version

	cliApp.Commands = []cli.Command{
		sync.Prepare(),
		app.Prepare(),
	}

	piaas.SubscribeExitSignal(func(sig os.Signal) {
		// Upon receiving the signal, wait for 500 millis before exit.
		time.AfterFunc(time.Millisecond*500, func() {
			log.Debugln("Piaas terminating ...")
			os.Exit(0)
		})
	}, false)
	piaas.HandleExitSignal()

	err := cliApp.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %s\n", err)
		os.Exit(1)
	}
}
