package app

import (
	"github.com/urfave/cli"
)

// Prepare will prepare the sync module.
// This will setup the relevant cli flags of this module.
func Prepare() cli.Command {
	return cli.Command{
		Name:    "app",
		Aliases: []string{"a"},
		Usage:   "Operating an app",
		Subcommands: []cli.Command{
			PrepareRun(),
			PrepareStop(),
		},
	}
}
