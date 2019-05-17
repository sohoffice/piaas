package piaas

import (
	"fmt"
	"github.com/sohoffice/piaas/util"
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

	expected := fmt.Sprintf("rsync -av %s --exclude-from=%s --delete --copy-links . %s", "-e 'ssh -o ConnectTimeout=10'", "/tmp/.piaasignore", "foo@bar:/tmp/foo")
	if tester.commands[0] != expected {
		t.Errorf("Unexpected sync all command.\nExpected: %s.\n  Actual: %s.", expected, tester.commands[0])
	}
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
	expected1 := fmt.Sprintf("rsync -av %s --exclude-from=%s %s --delete --copy-links . %s",
		"-e 'ssh -o ConnectTimeout=10'",
		"/tmp/.piaasignore",
		fmt.Sprintf("--include='*/' %s --exclude='*'", "--include='/a' --include='/b' --include='/c'"),
		"foo@bar:/tmp/foo")
	expected2 := fmt.Sprintf("rsync -av %s --exclude-from=%s %s --delete --copy-links . %s",
		"-e 'ssh -o ConnectTimeout=10'",
		"/tmp/.piaasignore",
		fmt.Sprintf("--include='*/' %s --exclude='*'", "--include='/d'"),
		"foo@bar:/tmp/foo")
	actual := util.StringArray(tester.commands)
	expected := util.StringArray([]string{expected1, expected2})

	if actual.Compare(expected) == false {
		t.Errorf("Unexpected sync files commands:\nExpected:\n%s\nActual:\n%s", expected.ToString(), actual.ToString())
	}
}

type RsyncWrapperTester struct {
	rs       RsyncWrapper
	commands []string
}

func prepare() *RsyncWrapperTester {
	rs := NewRsyncWrapper("rsync", "/tmp", "foo@bar:/tmp/foo")
	rs.SetIgnoreFile("/tmp/.piaasignore")

	tester := RsyncWrapperTester{
		rs: rs, commands: make([]string, 0),
	}

	// collect the commands
	tester.rs.Start(func(cmd string) {
		tester.commands = append(tester.commands, cmd)
	})

	return &tester
}
