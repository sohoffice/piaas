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

func main() {
	// flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	// Do not print any output from flag module.
	flag.CommandLine.SetOutput(ioutil.Discard)
	// We still would like the service of flag, for it will be used to initialize glog.
	flag.Parse()

	app := cli.NewApp()
	app.Name = "Piaas, tools to develop using multiple machines as if you have Personal IAAS."
	app.HelpName = "piaas"
	app.Authors = []cli.Author{
		{
			Name:  "Douglas Liu",
			Email: "douglas@sohoffice.com",
		},
	}
	app.Version = "v0.0.1"

	app.Commands = []cli.Command{
		sync.Prepare(),
	}

	defer glog.Flush()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
