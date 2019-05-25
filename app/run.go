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
	"time"
)

func PrepareRun() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r", "start"},
		Usage:   "Start the app",
		Flags: append(piaas.PrepareCommonFlags(), cli.BoolFlag{
			Name:  "tail, t",
			Usage: "Tail the app logs.",
		}),
		ArgsUsage: "[app name]",
		Action:    ExecuteRun,
	}
}

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
	runDir := piaas.NewRunDir("./.piaas.d")
	run(runDir, app, c)

	return nil
}

func run(runDir piaas.RunDir, app piaas.App, c *cli.Context) {
	var args []string
	for _, p := range app.Params {
		args = append(args, os.Expand(p, getEnvOrError))
	}

	logfileName := runDir.Logfile(app.Name)
	// delete existing logifle.
	os.Remove(logfileName)
	logfile, err := os.Create(logfileName)
	util.CheckError("create logfile", err)

	cmd := exec.Command(app.Cmd, args...)
	pid := runCmd(cmd, logfile)

	log.Infof("Run app %s: %s", app.Name, strings.Join(append([]string{app.Cmd}, args...), " "))
	log.Infof("Logfile: %s", logfileName)
	log.Infof("Pid: %d", pid)
	err = ioutil.WriteFile(runDir.Pidfile(app.Name), []byte(fmt.Sprintf("%d", pid)), 0644)
	util.CheckError("write pidfile", err)
	log.Infoln("---------------------------")

	if c.Bool("tail") {
		waitCh := make(chan bool, 1)

		tailer := tail(runDir, app)
		piaas.SubscribeExitSignal(func(sig os.Signal) {
			log.Debugf("Stop tailing ...")

			tailer.Close()
			waitCh <- true
		}, true)

		<-waitCh
		close(waitCh)
		log.Debugln("Stop piaas run.")
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

func tail(runDir piaas.RunDir, app piaas.App) *Tail {
	log.Infof("Tailing app %s: %s.", app.Name, runDir.Logfile(app.Name))
	tail := NewTail(runDir.Logfile(app.Name), os.Stdout, 10240)

	go func() {
		for {
			<-time.After(time.Millisecond * 500)
			tail.Read(-1)
		}
	}()

	tail.Start()
	tail.Read(-1)

	return tail
}
