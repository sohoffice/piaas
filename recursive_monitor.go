package piaas

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/util"
	"hash/fnv"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type MonitorPath struct {
	// the path relative to project root
	path string
}

type RecursiveMonitor struct {
	basedir string
	// This channel publish collected changes.
	// The changes will be published to CollectObservers.
	collectsCh chan []string

	// The caller of RecursiveMonitor should use this to subscribe to collected changes.
	collectObservers []chan<- []string

	// Monitors currently being watched
	monitors map[uint32]MonitorPath

	// Receive file system changes from fsnotify
	watcher *fsnotify.Watcher

	// This channel will publish the observed file changes.
	// Usually the channel receives file or directory name
	// However, it may also receive an empty string which means all accumulated changes should be collected
	changes chan string

	// The observers that subscribe to changes event.
	// Should only be used in testing
	changesObservers []chan<- string

	// The debouncer can be used to collect file change events over a small period of time.
	debouncer util.Debouncer

	// accumulated changes waiting to be collected.
	//
	// All values will be sent to collects channel when debounce event is triggered.
	// `accumulated` will be cleared thereafter.
	accumulated util.StringSet
}

// Start watching on all registered monitors, and manage the removal and addition of new directories.
func (rm *RecursiveMonitor) Start(debounceTime uint64) {
	rm.debouncer = util.NewDebouncer(debounceTime, func(tag string) {
		log.Debugf("Received debounce event: %s", tag)
		rm.changes <- ""
	})
	rm.collectsCh = make(chan []string)
	go collectedHandler(rm)
	go changesHandler(rm)
	go fsnotifyHandler(rm)
}

func (rm *RecursiveMonitor) SubscribeToChanges(subscriber chan<- string) {
	rm.changesObservers = append(rm.changesObservers, subscriber)
	log.Debugf("Added changes observer: %d.", len(rm.changesObservers))
}

func (rm *RecursiveMonitor) Subscribe(subscriber chan<- []string) {
	rm.collectObservers = append(rm.collectObservers, subscriber)
	log.Debugf("Added collect observer: %d.", len(rm.collectObservers))
}

func (rm *RecursiveMonitor) Stop() {
	rm.watcher.Close()
}

// add a directory to be monitored
func (rm *RecursiveMonitor) add(path string) {
	path = filepath.Clean(path)
	h := fnv.New32a()
	h.Write([]byte(path))
	hash := h.Sum32()
	m := MonitorPath{
		path: path,
	}
	rm.monitors[hash] = m

	err := rm.watcher.Add(path)
	if err != nil {
		log.Fatalf("Error adding path to RecursiveMonitor %s: %s", path, err)
	}
}

// Remove a potentially watched directory from the monitor.
// Return true if the directory was watched and removed.
func (rm *RecursiveMonitor) remove(path string) bool {
	path = filepath.Clean(path)
	h := fnv.New32a()
	h.Write([]byte(path))
	hash := h.Sum32()

	if _, ok := rm.monitors[hash]; ok {
		// if the path was monitored, remove it.
		err := rm.watcher.Remove(path)
		if err != nil {
			log.Debugf("Error removing path %s: %s", path, err)
		} else {
			delete(rm.monitors, hash)
			return true
		}
	} // do nothing if the path wasn't monitored.
	return false
}

// Number of watched directories
func (rm *RecursiveMonitor) length() int {
	return len(rm.monitors)
}

// Return the watched directories as a string array, in random order.
func (rm *RecursiveMonitor) watchedDirectories() []string {
	dir := make([]string, 0)
	for _, m := range rm.monitors {
		dir = append(dir, m.path)
	}
	return dir
}

// notify messages to all change observers
func (rm *RecursiveMonitor) notifyChanges(msg string) {
	// log.Debugf("notify %d changes observers: %s", len(rm.changesObservers), msg)
	for _, sub := range rm.changesObservers {
		sub <- msg
	}
}

func (rm *RecursiveMonitor) notifyCollected(collected []string) {
	log.Debugf("notify %d collect observers: %d messages.", len(rm.collectObservers), len(collected))
	for _, obs := range rm.collectObservers {
		obs <- collected
	}
}

// Listen to fsnotify events and post to the changes channel.
func fsnotifyHandler(rmPtr *RecursiveMonitor) {
	for {
		select {
		case event, ok := <-rmPtr.watcher.Events:
			if !ok {
				log.Debugf("Watcher event is not ok.")
				return
			}
			filename := filepath.Clean(event.Name)
			switch event.Op {
			case fsnotify.Create:
				info, err := os.Stat(filename)
				if err != nil {
					log.Infof("File error %s: %s", filename, err)
				} else if info.IsDir() { // a new directory was added
					rmPtr.add(filename)
				}
				rmPtr.changes <- filename
			case fsnotify.Remove:
				if rmPtr.remove(filename) {
					// this is actually a dead block, I don't believe we'll go into this in any situation.
					dir := filepath.Dir(filename)
					rmPtr.changes <- dir
				} else {
					rmPtr.changes <- filename
				}
			default:
				rmPtr.changes <- filename
			}
			log.Debugf("fsnotify event: %s", event)
		case err, ok := <-rmPtr.watcher.Errors:
			if !ok {
				log.Errorf("Watcher errors is not ok.")
				return
			}
			log.Errorf("error: %s", err)
		}
	}
}

func changesHandler(rmPtr *RecursiveMonitor) {
	for {
		msg := <-rmPtr.changes
		switch msg {
		case "": // collect event
			// log.Debugf("Collecting monitored changes: %s", rmPtr.accumulated)
			if rmPtr.accumulated == nil {
				// this is very wrong, accumulated should have at least one msg.
				log.Fatalln("Try to collect an empty accumulated list.")
			} else {
				rmPtr.collectsCh <- rmPtr.accumulated
				rmPtr.accumulated = nil
			}
		default:
			if rmPtr.accumulated == nil {
				rmPtr.accumulated = make([]string, 0)
			}
			if rmPtr.basedir != "" && !strings.HasPrefix(msg, rmPtr.basedir) {
				// make sure the path is full
				msg = path.Join(rmPtr.basedir, msg)
			}
			// rmPtr.accumulated = append(rmPtr.accumulated, msg)
			rmPtr.accumulated = *rmPtr.accumulated.Add(msg)
			rmPtr.debouncer.Event(msg)
			rmPtr.notifyChanges(msg)
		}
	}
}

// Relay the collected message to observers.
func collectedHandler(rmPtr *RecursiveMonitor) {
	for {
		collected := <-rmPtr.collectsCh
		rmPtr.notifyCollected(collected)
	}
}

// Traverse the tree and create monitor for all directory
func NewRecursiveMonitor(start string) RecursiveMonitor {
	monitors := make(map[uint32]MonitorPath)
	watcherPtr, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcherPtr: %s", err)
	}
	rm := RecursiveMonitor{
		basedir:          start,
		monitors:         monitors,
		watcher:          watcherPtr,
		changes:          make(chan string, 10),
		changesObservers: make([]chan<- string, 0),
	}

	log.Debugf("Walking: %s", start)
	err = filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Debugf("  | Added: %s", path)
			rm.add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the directory tree %s: %s.", start, err)
	}

	return rm
}
