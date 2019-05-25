// +build !windows

package app

import (
	"github.com/sohoffice/piaas/util"
	"os"
	"os/exec"
	"syscall"
)

func runCmd(cmd *exec.Cmd, output *os.File) int {
	// Use logfile for stdout and stderr
	cmd.Stdout = output
	cmd.Stderr = output

	// Run the app command in a separate process group to prevent it from being killed by ctrl-c
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	err := cmd.Start()
	util.CheckError("start command", err)

	return cmd.Process.Pid
}
