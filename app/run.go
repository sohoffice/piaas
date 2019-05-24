package app

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas"
	"github.com/sohoffice/piaas/util"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func PrepareRun() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r", "start"},
		Usage:   "Start the app",
		Flags: append(piaas.PrepareCommonFlags(), cli.BoolFlag{
			Name:  "tail, t",
			Usage: "Tail the output after the app is started.",
		}),
		ArgsUsage: "[app name]",
		Action:    ExecuteRun,
	}
}

var runDir piaas.RunDir

func ExecuteRun(c *cli.Context) error {
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
	runDir = piaas.NewRunDir("./.piaas.d")
	run(app, c)

	return nil
}

func run(app piaas.App, c *cli.Context) {
	var params []string
	for _, p := range app.Params {
		params = append(params, os.Expand(p, getEnvOrError))
	}

	logfileName := runDir.Logfile(app.Name)
	// delete existing logifle.
	os.Remove(logfileName)
	logfile, err := os.Create(logfileName)
	util.CheckError("create logfile", err)

	cmd := exec.Command(app.Cmd, params...)
	// Use logfile for stdout and stderr
	cmd.Stdout = logfile
	cmd.Stderr = logfile

	err = cmd.Start()
	util.CheckError("start command", err)
	pid := fmt.Sprintf("%d", cmd.Process.Pid)

	log.Infof("Run app %s: %s", app.Name, strings.Join(append([]string{app.Cmd}, params...), " "))
	log.Infof("Logfile: %s", logfileName)
	log.Infof("Pid: %d", cmd.Process.Pid)
	err = ioutil.WriteFile(runDir.Pidfile(app.Name), []byte(pid), 0644)
	util.CheckError("write pidfile", err)

	if c.Bool("tail") {
		// Do not complete in tail mode.
		completeCh := make(chan bool, 1)

		// wait for a channel that will never arrive.
		<-completeCh
	}
}

func getEnvOrError(name string) string {
	val, exists := os.LookupEnv(name)
	if !exists {
		log.Errorf("Environment variable not found: %s", name)
		val = ""
	}
	return val
}
