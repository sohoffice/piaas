package piaas

import (
	"fmt"
	"github.com/sohoffice/piaas/stringarrays"
	"os/exec"
	"testing"
	"time"
)

func TestRsyncWrapper_SyncAll(t *testing.T) {
	tester := prepare()

	tester.rs.SyncAll()
	// wait a small while to make sure channel has received the message
	<-time.After(time.Millisecond * 200)

	if len(tester.commands) != 1 {
		t.Errorf("Unexpected command length, should be 1, but is %d.", len(tester.commands))
	}

	args := []string{"-av", "-e", "ssh -o ConnectTimeout=10", fmt.Sprintf("--exclude-from=%s", "/tmp/.piaasignore"), "--delete", "--copy-links", ".", "foo@bar:/tmp/foo"}
	expected := exec.Command("rsync", args...)
	actual := tester.commands[0]
	compareCmd(t, expected, actual)
}

func TestRsyncWrapper_SyncFiles(t *testing.T) {
	tester := prepare()

	tester.rs.SyncFiles([]string{})
	<-time.After(time.Millisecond * 50)
	if len(tester.commands) != 0 {
		t.Errorf("Commands should not have been collected.")
	}

	tester.rs.SyncFiles([]string{"/a", "/b", "/c"})
	tester.rs.SyncFiles([]string{"/d"})
	<-time.After(time.Millisecond * 200)

	if len(tester.commands) != 2 {
		t.Errorf("Unexpected command length, should be 2, but is %d.", len(tester.commands))
	}
	args1 := []string{"-av", "-e", "ssh -o ConnectTimeout=10", fmt.Sprintf("--exclude-from=%s", "/tmp/.piaasignore"),
		"--include='*/'", "--include='/a'", "--include='/b'", "--include='/c'", "--exclude='*'",
		"--delete", "--copy-links", ".", "foo@bar:/tmp/foo"}
	expected1 := exec.Command("rsync", args1...)

	args2 := []string{"-av", "-e", "ssh -o ConnectTimeout=10", fmt.Sprintf("--exclude-from=%s", "/tmp/.piaasignore"),
		"--include='*/'", "--include='/d'", "--exclude='*'",
		"--delete", "--copy-links", ".", "foo@bar:/tmp/foo"}
	expected2 := exec.Command("rsync", args2...)
	actuals := tester.commands
	expecteds := []*exec.Cmd{expected1, expected2}

	for i := range expecteds {
		actual := actuals[i]
		expected := expecteds[i]
		compareCmd(t, expected, actual)
	}
}

type RsyncWrapperTester struct {
	rs       RsyncWrapper
	commands []*exec.Cmd
}

func prepare() *RsyncWrapperTester {
	rs := NewRsyncWrapper("rsync", "/tmp", "foo@bar:/tmp/foo")
	rs.SetIgnoreFile("/tmp/.piaasignore")

	tester := RsyncWrapperTester{
		rs: rs, commands: make([]*exec.Cmd, 0),
	}

	// collect the commands
	tester.rs.Start(func(cmd *exec.Cmd) {
		tester.commands = append(tester.commands, cmd)
	})

	return &tester
}

func compareCmd(t *testing.T, expected *exec.Cmd, actual *exec.Cmd) {
	if expected.Path != actual.Path {
		t.Errorf("Unexpected command path.\nExpected: %s.\n  Actual: %s.", expected.Path, actual.Path)
	}
	if stringarrays.Compare(expected.Args, actual.Args) == false {
		t.Errorf("Unexpected command args.\nExpected: %s.\n  Actual: %s.", stringarrays.ToString(expected.Args), stringarrays.ToString(actual.Args))
	}
}
