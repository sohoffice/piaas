package piaas

import (
	"github.com/sohoffice/piaas/util"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRunDir(t *testing.T) {
	tempDir := os.TempDir()
	// We'd like RunDir to create the directory, but we wanted to use the temp path.
	os.RemoveAll(tempDir)
	defer os.RemoveAll(tempDir)

	dir := NewRunDir(tempDir)
	checkPath(t, path.Join(tempDir, "foo", LogfileName), dir.Logfile("foo"))
	checkPath(t, path.Join(tempDir, "foo", PidfileName), dir.Pidfile("foo"))
}

func checkPath(t *testing.T, expected string, actual string) {
	actual2, err := filepath.EvalSymlinks(actual)
	util.CheckError("check actual path", err)
	expected2, err := filepath.EvalSymlinks(expected)
	util.CheckError("check expected path", err)

	if expected2 != actual2 {
		t.Errorf("Unexpected path. Expecting %s, Actual: %s.", expected, actual)
	}
}
