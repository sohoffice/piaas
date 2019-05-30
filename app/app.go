package app

import (
	"fmt"
	"github.com/sohoffice/piaas"
	"github.com/urfave/cli"
)

var defaultRunDir = piaas.NewRunDir("./.piaas.d")

// Prepare will prepare the sync module.
// This will setup the relevant cli flags of this module.
func Prepare() cli.Command {
	return cli.Command{
		Name:    "app",
		Aliases: []string{"a"},
		Usage:   "Operating an app",
		Flags:   piaas.PrepareCommonFlags(),
		Subcommands: []cli.Command{
			PrepareRun(),
			PrepareStop(),
			PrepareStatus(),
		},
		Action: ExecuteApp,
	}
}

func ExecuteApp(c *cli.Context) error {
	if c.NArg() > 0 {
		fmt.Println("Invalid app command")
		cli.ShowAppHelpAndExit(c, 1)
	}

	err := piaas.HandleDebug(c)
	if err != nil {
		return err
	}

	config := piaas.ReadConfig(c.String("config"))
	printAllStatus(defaultRunDir, config)

	return nil
}
