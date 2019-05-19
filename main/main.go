package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/sohoffice/piaas/sync"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Piaas, tools to develop with multiple machines as if you have Personal IAAS."
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
	var debug bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Print debug messages",
			Destination: &debug,
		},
	}
	if debug == true {
		flag.Parse()
		flag.Lookup("logtostderr").Value.Set("true")
	}
	defer glog.Flush()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
