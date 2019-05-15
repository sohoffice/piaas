package piaas

import (
	"github.com/sohoffice/piaas/util"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

func TestNewRecursiveMonitor(t *testing.T) {
	var mt = MonitorTest(*t)
	tempDir := mt.prepareTestDir()
	defer os.RemoveAll(tempDir)

	expectedMonitorNames := util.StringArray([]string{
		tempDir, path.Join(tempDir, "foo"), path.Join(tempDir, "foo", "bar"),
		path.Join(tempDir, "foo", "baz"), path.Join(tempDir, "foo1"), path.Join(tempDir, "to-be-deleted"),
		path.Join(tempDir, "to-be-renamed-dir"),
	})
	monitors := NewRecursiveMonitor(tempDir)
	if monitors.length() != len(expectedMonitorNames) {
		t.Errorf("monitor number should be %d, but is %d.\nExpected:\n%s\nActual:\n%s", len(expectedMonitorNames), monitors.length(),
			expectedMonitorNames.ToString(), util.StringArray(monitors.watchedDirectories()).ToString())
	}

	for _, monitoredPath := range monitors.watchedDirectories() {
		if expectedMonitorNames.IndexOf(monitoredPath) == -1 {
			t.Errorf("monitor path %s was not expected.", monitoredPath)
		}
	}
}

// Validate the file changes were actually captured by the recursive monitor.
func TestMonitorFileChanges(t *testing.T) {
	var mt = MonitorTest(*t)
	mtPtr := &mt
	tempDir := mt.prepareTestDir()
	defer os.RemoveAll(tempDir)

	ch := make(chan bool)
	subscribe := make(chan string)
	expectedChanges := util.StringSet(make([]string, 0))
	rm := NewRecursiveMonitor(tempDir)
	// Setup changes observer
	changes := util.StringSet(make([]string, 0))
	go func() {
		for {
			msg := <-subscribe
			log.Printf("Observed change: %s.", msg)
			changes = *changes.Add(msg)
		}
	}()

	// add myself to the distribution list of changes event
	rm.SubscribeToChanges(subscribe)
	// start watching
	rm.Watch()

	// making changes
	// creating new files
	mtPtr.touchFile(path.Join(tempDir, "foo", "foo-file"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo", "foo-file"))
	mtPtr.touchFile(path.Join(tempDir, "baz-file"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "baz-file"))
	// update file
	bytes := []byte("foo")
	mtPtr.writeFile(path.Join(tempDir, "foo_file"), &bytes)
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo_file"))
	// delete a file
	mtPtr.removeFile(path.Join(tempDir, "foo", "abc"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo", "abc"))
	// delete a directory
	mtPtr.removeFile(path.Join(tempDir, "to-be-deleted"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "to-be-deleted"))
	// rename a directory
	mtPtr.rename(path.Join(tempDir, "to-be-renamed-dir"), path.Join(tempDir, "renamed"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "to-be-renamed-dir"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "renamed"))
	// rename a file
	mtPtr.rename(path.Join(tempDir, "foo", "to-be-renamed"), path.Join(tempDir, "foo", "renamed"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo", "to-be-renamed"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo", "renamed"))
	// test by adding a new directory
	mtPtr.mkDir(path.Join(tempDir, "foo-dir"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo-dir"))
	<-time.After(time.Millisecond * 50)
	// the below file can be added before foo-dir was monitored, so we have to wait a few while to make sure the subscription works.
	mtPtr.touchFile(path.Join(tempDir, "foo-dir", "abc"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo-dir", "abc"))
	go func() {
		// stop the test after 500 millis
		<-time.After(time.Millisecond * 500)
		defer func() {
			ch <- true
		}()

		log.Printf("Event validating observed changes: %s", changes)
		if len(changes) != len(expectedChanges) {
			t.Fatalf("Should receive %d changes, but %d was received.\n%s", len(expectedChanges), len(changes), changes)
		}
		if !expectedChanges.Compare(changes) {
			t.Fatalf("Expecting: \n%s\nActual: \n%s", expectedChanges.ToString(), changes.ToString())
		}

		log.Print("\n\n\n")
	}()
	<-ch
}

type MonitorTest testing.T

func (t *MonitorTest) mkDir(path string) {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		t.Fatalf("Error creating test directory %s: %s", path, err)
	}
}

func (t *MonitorTest) touchFile(path string) {
	var bytes []byte
	t.writeFile(path, &bytes)
}

// write something to file
func (t *MonitorTest) writeFile(path string, bytes *[]byte) {
	err := ioutil.WriteFile(path, *bytes, 0644)
	if err != nil {
		t.Fatalf("Error touching file %s: %s", path, err)
	}
}

func (t *MonitorTest) removeFile(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Error deleting file %s: %s", path, err)
	}
}

func (t *MonitorTest) rename(oldName string, newName string) {
	err := os.Rename(oldName, newName)
	if err != nil {
		t.Fatalf("Error renaming file %s to %s: %s", oldName, newName, err)
	}
}

func (t *MonitorTest) prepareTestDir() string {
	tempDir, err := ioutil.TempDir("", "walk-test")
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}

	// Create the test tree hierarchy
	t.mkDir(path.Join(tempDir, "foo", "bar"))
	t.mkDir(path.Join(tempDir, "foo", "baz"))
	t.touchFile(path.Join(tempDir, "foo", "to-be-renamed"))
	t.mkDir(path.Join(tempDir, "to-be-deleted"))
	t.mkDir(path.Join(tempDir, "to-be-renamed-dir"))
	t.mkDir(path.Join(tempDir, "foo1"))
	t.touchFile(path.Join(tempDir, "foo_file"))
	t.touchFile(path.Join(tempDir, "foo", "abc"))

	return tempDir
}
