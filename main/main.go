package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/sohoffice/piaas/sync"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

var version string

func main() {
	// Override flag.CommandLine to prevent extraneous flag from stopping the execution. It should be managed by urfave/cli.
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	// Do not print any output from flag module.
	flag.CommandLine.SetOutput(ioutil.Discard)
	// We still would like the service of flag, for it will be used to initialize glog.
	flag.Parse()

	app := cli.NewApp()
	app.Name = "Piaas, tools to develop using multiple machines as if using Personal IAAS."
	app.HelpName = "piaas"
	app.Authors = []cli.Author{
		{
			Name:  "Douglas Liu",
			Email: "douglas@sohoffice.com",
		},
	}
	app.Version = version

	app.Commands = []cli.Command{
		sync.Prepare(),
	}

	defer glog.Flush()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %s\n", err)
		os.Exit(1)
	}
}
