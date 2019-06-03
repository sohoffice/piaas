// +build windows

package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/util"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func runCmd(cmd *exec.Cmd, output *os.File) int {
	// Use logfile for stdout and stderr
	cmd.Stdout = output
	cmd.Stderr = output

	// Run the app command in a separate process group to prevent it from being killed by ctrl-c
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	log.Debugf("Run windows command: %s", strings.Join(cmd.Args, " "))
	err := cmd.Start()
	util.CheckError("start command", err)

	return cmd.Process.Pid
}
