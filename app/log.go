package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas"
	"github.com/urfave/cli"
	"os"
)

// PrepareLog prepare command line argument for the log sub command
func PrepareLog() cli.Command {
	return cli.Command{
		Name:      "log",
		Aliases:   []string{"log"},
		Usage:     "Tail the app logs",
		ArgsUsage: "[app name]",
		Flags:     piaas.PrepareCommonFlags(),
		Action:    ExecuteLog,
	}
}

// ExecuteLog is the entry point of the log sub command.
func ExecuteLog(c *cli.Context) error {
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

	return logApp(defaultRunDir, app)
}

func logApp(runDir piaas.RunDir, app piaas.App) error {
	_, _, err := findAppProc(runDir, app)
	if err != nil {
		// remove the pidfile if the process can not be located.
		defer os.Remove(runDir.Pidfile(app.Name))
		return fmt.Errorf("app '%s' is not running", app.Name)
	}

	waitCh := make(chan bool, 1)

	tailer := tail(runDir, app)
	piaas.SubscribeExitSignal(func(sig os.Signal) {
		log.Debugf("Stop tailing ...")

		tailer.Close()
		waitCh <- true
	}, true)

	<-waitCh
	close(waitCh)
	log.Debugln("Stop piaas log.")

	return nil
}
