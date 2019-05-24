package app

import (
	"github.com/urfave/cli"
)

// Prepare the sync module.
// This usually involves setup the correct module flags.
func Prepare() cli.Command {
	return cli.Command{
		Name:    "app",
		Aliases: []string{"a"},
		Usage:   "Operating an app",
		Subcommands: []cli.Command{
			PrepareRun(),
		},
	}
}
