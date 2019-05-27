package app

import (
	"bufio"
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tempDir := os.TempDir()
	os.MkdirAll(tempDir, 0755)
	log.Debugln(tempDir)
	// defer os.RemoveAll(tempDir)

	runDir := piaas.NewRunDir(tempDir)
	var app1 = piaas.App{
		Name: "sleep",
	}
	logfile1 := runDir.Logfile("sleep")
	os.Remove(logfile1)

	if runtime.GOOS == "windows" {
		// app1.Cmd = "cmd.exe"
		// app1.Params = []string{"/c", "echo sleeping & timout /t 1"}
		batch, _ := filepath.Abs(path.Join(".", "wintest.bat"))
		app1.Cmd = "cmd.exe"
		app1.Params = []string{"/c", batch}
		log.Infof("Testing on windows, use %s %s", app1.Cmd, strings.Join(app1.Params, " "))
	} else {
		app1.Cmd = "bash"
		app1.Params = []string{"-c", "echo sleeping && sleep 1"}
	}

	cliApp := cli.NewApp()
	ctx := cli.NewContext(cliApp, flag.NewFlagSet("test", flag.ContinueOnError), nil)

	run(runDir, app1, ctx)

	pidfile := runDir.Pidfile("sleep")
	pid, _ := ioutil.ReadFile(pidfile)
	pidInt, _ := strconv.ParseInt(string(pid), 10, 32)
	proc, err := os.FindProcess(int(pidInt))
	if err != nil {
		t.Errorf("Error finding process, pid: %s.: %s", pid, err)
	}
	// wait a small while to make sure the logfile is written
	time.Sleep(time.Millisecond * 300)
	proc.Kill()

	logData1, err := ioutil.ReadFile(logfile1)
	if err != nil {
		t.Errorf("Error reading logfile content, %s", err)
	}
	logStr := string(logData1)
	scanner := bufio.NewScanner(strings.NewReader(logStr))
	if !scanner.Scan() {
		t.Errorf("The file is empty.")
	} else {
		firstLine := scanner.Text()
		if firstLine != "sleeping" {
			t.Errorf("Unexpected log file content. first line: %s, all: %s.", firstLine, logStr)
		}
	}

	// Test environment variable substitution
	logfile2 := runDir.Logfile("echo")
	os.Remove(logfile2)
	app2 := piaas.App{
		Name: "echo",
	}
	greeting := "Hi, ${SALUTE}${TEST_GREETING}. Today is a good day to code."
	if runtime.GOOS == "windows" {
		app2.Cmd = "cmd.exe"
		app2.Params = []string{"/c", "echo " + greeting}
	} else {
		app2.Cmd = "echo"
		app2.Params = []string{greeting}
	}
	os.Setenv("TEST_GREETING", "tester")
	run(runDir, app2, ctx)
	// wait a small while to make sure the logfile is written
	time.Sleep(time.Millisecond * 300)

	logData2, err := ioutil.ReadFile(logfile2)
	if err != nil {
		t.Errorf("Error reading logfile content, %s", err)
	}
	scanner2 := bufio.NewScanner(strings.NewReader(string(logData2)))
	if !scanner2.Scan() {
		t.Error("The file is empty.")
	} else {
		firstLine2 := scanner2.Text()
		if firstLine2 != "Hi, tester. Today is a good day to code." {
			t.Errorf("Unexpected log file content: %s.", firstLine2)
		}
	}
}
