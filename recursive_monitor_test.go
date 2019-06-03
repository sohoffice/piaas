package piaas

import (
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/stringarrays"
	"github.com/sohoffice/piaas/util"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"testing"
	"time"
)

func TestNewRecursiveMonitor(t *testing.T) {
	var mt = MonitorTest(*t)
	tempDir := mt.prepareTestDir()
	defer os.RemoveAll(tempDir)

	expectedMonitorNames := []string{
		tempDir, path.Join(tempDir, "foo"), path.Join(tempDir, "foo", "bar"),
		path.Join(tempDir, "foo", "baz"), path.Join(tempDir, "foo1"), path.Join(tempDir, "to-be-deleted"),
		path.Join(tempDir, "to-be-renamed-dir"),
	}
	rm := NewRecursiveMonitor(tempDir)
	if rm.length() != len(expectedMonitorNames) {
		t.Errorf("monitor number should be %d, but is %d.\nExpected:\n%s\nActual:\n%s", len(expectedMonitorNames), rm.length(),
			stringarrays.ToString(expectedMonitorNames), stringarrays.ToString(rm.watchedDirectories()))
	}

	actual := util.StringSet{}
	for _, act := range rm.watchedDirectories() {
		log.Debugf("  | Add watched: %s", filepath.ToSlash(act))
		actual = *actual.Add(filepath.ToSlash(act))
	}
	for _, expected := range expectedMonitorNames {
		exp := filepath.ToSlash(expected)
		if stringarrays.IndexOf(actual, exp) == -1 {
			t.Errorf("Expected monitor path %s was not found. Actual: %s", exp, stringarrays.ToString(actual))
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
			// log.Printf("Observed change: %s.", msg)
			changes = *changes.Add(msg)
		}
	}()

	// add myself to the distribution list of changes event
	rm.SubscribeToChanges(subscribe)
	// start watching
	rm.Start(1000)

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
	<-time.After(time.Millisecond * 200)
	// the below file can be added before foo-dir was monitored, so we have to wait a few while to make sure the subscription works.
	mtPtr.touchFile(path.Join(tempDir, "foo-dir", "abc"))
	expectedChanges = *expectedChanges.Add(path.Join(tempDir, "foo-dir", "abc"))
	go func() {
		// stop the test after 500 millis
		<-time.After(time.Millisecond * 500)
		defer func() {
			ch <- true
		}()

		log.Debugf("Event validating observed changes: %s", changes)
		if len(changes) < len(expectedChanges) {
			t.Fatalf("Should receive %d changes, but %d was received.\n%s", len(expectedChanges), len(changes), changes)
		}
		for _, exp := range expectedChanges {
			expected := filepath.FromSlash(exp)
			if changes.IndexOf(expected) == -1 {
				t.Errorf("Expected change %s wasn't recorded.", expected)
			}
		}
	}()
	<-ch
}

// Make sure RecursiveMonitor implements the Monitor interface.
func TestMonitorInterface(t *testing.T) {
	var mt = MonitorTest(*t)
	tempDir := mt.prepareTestDir()
	defer os.RemoveAll(tempDir)

	rm := NewRecursiveMonitor(tempDir)
	var monitor Monitor = &rm
	ch := make(chan []string)
	go func() {
		collected := <-ch
		var converted []string
		for _, col := range collected {
			converted = append(converted, filepath.ToSlash(col))
		}
		sort.Strings(converted)

		if stringarrays.IndexOf(converted, filepath.ToSlash(path.Join(tempDir, "collect1"))) == -1 ||
			stringarrays.IndexOf(converted, filepath.ToSlash(path.Join(tempDir, "collect2"))) == -1 ||
			stringarrays.IndexOf(converted, filepath.ToSlash(path.Join(tempDir, "collect3"))) == -1 {
			log.Errorln("Collected results:")
			log.Errorln(stringarrays.ToString(converted))
			t.Errorf("Unexpected collect results: %s", stringarrays.ToString(converted))
		}
	}()
	monitor.Subscribe(ch)
	monitor.Start(300)

	mt.touchFile(path.Join(tempDir, "collect1"))
	mt.touchFile(path.Join(tempDir, "collect2"))
	mt.touchFile(path.Join(tempDir, "collect3"))

	// wait for a small while to make sure collected are checked.
	<-time.After(time.Millisecond * 500)
	monitor.Stop()
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
		log.Errorf("Error touching file %s: %s", path, err)
	}
}

func (t *MonitorTest) removeFile(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Errorf("Error deleting file %s: %s", path, err)
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
		log.Errorf("Error creating temp dir: %s", err)
	}

	// Create the test tree hierarchy
	t.mkDir(path.Join(tempDir, "foo", "bar"))
	t.mkDir(path.Join(tempDir, "foo", "baz"))
	t.mkDir(path.Join(tempDir, "foo1"))
	t.mkDir(path.Join(tempDir, "to-be-deleted"))
	t.mkDir(path.Join(tempDir, "to-be-renamed-dir"))
	t.touchFile(path.Join(tempDir, "foo", "to-be-renamed"))
	t.touchFile(path.Join(tempDir, "foo_file"))
	t.touchFile(path.Join(tempDir, "foo", "abc"))

	return tempDir
}
