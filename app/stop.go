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

// PrepareStop prepare command line argument for stop sub command
func PrepareStop() cli.Command {
	return cli.Command{
		Name:      "stop",
		Aliases:   []string{"s"},
		Usage:     "Stop the app",
		ArgsUsage: "[app name]",
		Flags:     piaas.PrepareCommonFlags(),
		Action:    ExecuteStop,
	}
}

// ExecuteStop is the entry point of stop sub command.
func ExecuteStop(c *cli.Context) error {
	var appName string
	if c.NArg() >= 1 {
		appName = c.Args().Get(0)
	}
	err := piaas.HandleDebug(c)
	if err != nil {
		return err
	}

	config := piaas.ReadConfig(c.String("config"))
	app, err := config.GetApp(appName)
	if err != nil {
		return err
	}

	return stop(defaultRunDir, app)
}

// stop actually stop the process noted in the pidfile.
func stop(runDir piaas.RunDir, app piaas.App) error {
	pidBytes, err := ioutil.ReadFile(runDir.Pidfile(app.Name))
	if err != nil {
		log.Debugf("Error: %s", err)
		return fmt.Errorf("app '%s' is not running", app.Name)
	}

	pid, err := strconv.ParseInt(string(pidBytes), 10, 32)
	if err != nil {
		log.Debugf("Error: %s", err)
		return fmt.Errorf("app '%s' is not running", app.Name)
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		log.Debugf("Error: %s", err)
		return fmt.Errorf("app '%s' is not running", app.Name)
	}

	err = proc.Kill()
	if err != nil {
		log.Debugf("Error: %s", err)
		return fmt.Errorf("can not kill app '%s'", app.Name)
	}

	// remove the pidfile on successful process kill.
	defer os.Remove(runDir.Pidfile(app.Name))
	log.Infof("App '%s' stopped.", app.Name)

	return nil
}
