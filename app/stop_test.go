package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
)

func TestStop(t *testing.T) {
	tempDir := os.TempDir()
	defer os.RemoveAll(tempDir)
	runDir := piaas.NewRunDir(tempDir)
	// Use a very large pid, this should less likely to be used by modern operating system
	ioutil.WriteFile(runDir.Pidfile("stop1"), []byte("9999999"), 0644)

	err := stop(runDir, piaas.App{
		Name: "stop1",
	})
	if err == nil {
		t.Errorf("Should return error ")
	} else {
		log.Debugf("Expected error: %s", err)
	}
	_, err = os.Stat(runDir.Pidfile("stop1"))
	if err != nil {
		t.Errorf("Pidfile has error: %s", err)
	} else {
		log.Debugln("Pidfile remain to exist if fail to stop the app.")
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		batch, _ := filepath.Abs(path.Join(".", "wintest.bat"))
		cmd = exec.Command("cmd.exe", "/c", batch)
	} else {
		cmd = exec.Command("sleep", "1")
	}
	cmd.Start()
	pid := strconv.FormatInt(int64(cmd.Process.Pid), 10)
	ioutil.WriteFile(runDir.Pidfile("stop2"), []byte(pid), 0644)
	err = stop(runDir, piaas.App{
		Name: "stop2",
	})
	if err != nil {
		t.Errorf("Fail to stop the process: %s", err)
	} else {
		log.Debugln("Successfully stopped the command")
	}
	_, err = os.Stat(runDir.Pidfile("stop2"))
	if err != nil {
		log.Debugln("Pidfile no longer exist if app is stopped:", err)
	} else {
		t.Errorf("Pidfile has error: %s", err)
	}
}
