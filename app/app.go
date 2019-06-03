package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"strconv"
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
			PrepareLog(),
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

// findAppProc finds the app process by reading the pidfile.
func findAppProc(runDir piaas.RunDir, app piaas.App) (*os.Process, int, error) {
	pidBytes, err := ioutil.ReadFile(runDir.Pidfile(app.Name))
	if err != nil {
		log.Debugf("Error reading pidfile: %s", err)
		return nil, -1, err
	}

	pid, err := strconv.ParseInt(string(pidBytes), 10, 32)
	if err != nil {
		log.Debugf("Error parsing pid: %s", err)
		return nil, -1, err
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		log.Debugf("Error locating process: %s", err)
		return nil, -1, err
	}

	return proc, int(pid), nil
}
