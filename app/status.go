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

type status uint8

const (
	Running status = iota
	Stopped
)

// PrepareStatus will print status of known apps, or the specified app if the name is specified.
func PrepareStatus() cli.Command {
	return cli.Command{
		Name:      "status",
		Usage:     "Print the app status",
		ArgsUsage: "[app names ...]",
		Flags:     piaas.PrepareCommonFlags(),
		Action:    ExecuteStatus,
	}
}

// ExecuteStop is the entry point of stop sub command.
func ExecuteStatus(c *cli.Context) error {
	err := piaas.HandleDebug(c)
	if err != nil {
		return err
	}

	config := piaas.ReadConfig(c.String("config"))

	if c.NArg() >= 1 {
		for i := 0; i < c.NArg(); i++ {
			appName := c.Args().Get(i)
			app, err := config.GetApp(appName)
			if err != nil {
				log.Debugln("Error reading app:", err)
				return fmt.Errorf("error finding app '%s'", appName)
			}
			printStatus(defaultRunDir, app)
		}
	} else {
		printAllStatus(defaultRunDir, config)
	}

	return nil
}

// printAllStatus will print status of all defined apps
func printAllStatus(runDir piaas.RunDir, config piaas.Config) {
	fmt.Printf("%10s: %s\n", "App", "Status")
	fmt.Println("------------------------------")
	for _, app := range config.Apps {
		printStatus(runDir, app)
	}
}

// printStatus print status of the specified app
func printStatus(dir piaas.RunDir, app piaas.App) {
	pidFilename := dir.Pidfile(app.Name)
	pidBytes, err := ioutil.ReadFile(pidFilename)
	if err != nil {
		doPrintStatus(app, Stopped, -1)
		return
	}
	pid, err := strconv.ParseInt(string(pidBytes), 10, 32)
	if err != nil {
		doPrintStatus(app, Stopped, -1)
		return
	}
	_, err = os.FindProcess(int(pid))
	if err != nil {
		doPrintStatus(app, Stopped, -1)
		return
	}
	doPrintStatus(app, Running, int(pid))
}

// doPrintStatus actually format and print the status string of an app
func doPrintStatus(app piaas.App, status status, pid int) {
	if status == Running {
		fmt.Printf("%10s: %s, pid %d\n", app.Name, "Running", pid)
	} else {
		fmt.Printf("%10s: %s\n", app.Name, "Stopped")
	}
}
