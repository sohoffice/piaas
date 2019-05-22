// +build darwin

package piaas

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/stringarrays"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestFSEventMonitor(t *testing.T) {
	tempDir := prepareFSeventMonitorTestDir(t)
	defer os.RemoveAll(tempDir)

	completed := make(chan bool)
	fsm := NewFSEventMonitor(tempDir)
	var mon Monitor = &fsm
	collected := make([]string, 0)
	collectedCh := make(chan []string)
	go func() {
		events := <-collectedCh
		collected = append(collected, events...)

		// start validating the test results
		completed <- true
	}()

	mon.Subscribe(collectedCh)
	mon.Start(600)

	expected := []string{
		pathConvert(path.Join(tempDir, "dir-to-delete")),
		pathConvert(path.Join(tempDir, "foo", "file-to-delete")),
	}
	// create a new directory
	mkDir(t, path.Join(tempDir, "dir"))
	// add a file in the newly created directory
	touchFile(t, path.Join(tempDir, "dir", "foo"))
	// delete a directory
	removeFile(t, path.Join(tempDir, "dir-to-delete"))
	// delete a file
	removeFile(t, path.Join(tempDir, "foo", "file-to-delete"))

	<-completed

	log.Debugf("Collected FSEvents:\n%s\n", stringarrays.ToString(collected))
	expected = append(expected, path.Join(tempDir, "dir"), path.Join(tempDir, "dir", "foo"))
	validate(t, collected, expected, []string{}, pathConvert)
}

func validate(t *testing.T, all []string, positives []string, negatives []string, conv func(string) string) {
	errFlag := false
	for _, s := range positives {
		s = conv(s)
		if stringarrays.IndexOf(all, s) == -1 {
			errFlag = true
			t.Errorf("Expected %s to exists, but was not.", s)
		}
	}
	for _, s := range negatives {
		s = conv(s)
		if stringarrays.IndexOf(all, s) != -1 {
			errFlag = true
			t.Errorf("Unexpected element %s.", s)
		}
	}
	if errFlag {
		fmt.Fprintf(os.Stderr, stringarrays.ToString(all))
	}
}

func pathConvert(path string) string {
	evaluated, err := filepath.EvalSymlinks(path)
	if err != nil {
		evaluated = path
	}
	return evaluated
}

func mkDir(t *testing.T, path string) {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		t.Fatalf("Error creating test directory %s: %s", path, err)
	}
}

func touchFile(t *testing.T, path string) {
	var bytes []byte
	writeFile(t, path, &bytes)
}

// write something to file
func writeFile(t *testing.T, path string, bytes *[]byte) {
	err := ioutil.WriteFile(path, *bytes, 0644)
	if err != nil {
		t.Fatalf("Error touching file %s: %s", path, err)
	}
}

func removeFile(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Error deleting file %s: %s", path, err)
	}
}

func rename(t *testing.T, oldName string, newName string) {
	err := os.Rename(oldName, newName)
	if err != nil {
		t.Fatalf("Error renaming file %s to %s: %s", oldName, newName, err)
	}
}

func prepareFSeventMonitorTestDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "fsevents")
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}

	// Create the test tree hierarchy
	mkDir(t, path.Join(tempDir, "foo", "bar"))
	mkDir(t, path.Join(tempDir, "foo", "baz"))
	touchFile(t, path.Join(tempDir, "foo", "file-to-delete"))
	touchFile(t, path.Join(tempDir, "foo", "file-to-rename"))
	mkDir(t, path.Join(tempDir, "dir-to-delete"))
	mkDir(t, path.Join(tempDir, "dir-to-rename"))
	mkDir(t, path.Join(tempDir, "foo1"))
	touchFile(t, path.Join(tempDir, "foo_file"))
	touchFile(t, path.Join(tempDir, "foo", "abc"))

	return tempDir
}
