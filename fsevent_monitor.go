// +build darwin

package piaas

import (
	"github.com/fsnotify/fsevents"
	"github.com/golang/glog"
	"log"
	"path/filepath"
	"time"
)

type FSEventMonitor struct {
	startDir string

	// This channel publish collected changes.
	// The changes will be published to CollectObservers.
	collectsCh chan []string

	// The caller of RecursiveMonitor should use this to subscribe to collected changes.
	collectObservers []chan<- []string

	eventStream *fsevents.EventStream
}

func (fsm *FSEventMonitor) Start(debounceTime uint64) {
	dev, err := fsevents.DeviceForPath(fsm.startDir)
	if err != nil {
		log.Fatal(err)
	}

	es := &fsevents.EventStream{
		Paths:   []string{fsm.startDir},
		Latency: time.Millisecond * time.Duration(debounceTime),
		Device:  dev,
		Flags:   fsevents.FileEvents | fsevents.NoDefer,
	}

	fsm.eventStream = es

	go handleCollectedFSEvent(fsm)
	go handleFSEvent(fsm)

	es.Start()
}

func (fsm *FSEventMonitor) Subscribe(subscriber chan<- []string) {
	fsm.collectObservers = append(fsm.collectObservers, subscriber)
	glog.Infof("Added collect observer: %d.", len(fsm.collectObservers))
}

func (fsm *FSEventMonitor) Stop() {
	if fsm.eventStream != nil {
		fsm.eventStream.Stop()
	}
}

func handleFSEvent(fsm *FSEventMonitor) {
	for events := range fsm.eventStream.Events {
		collected := make([]string, len(events))
		for i := range events {
			collected[i] = "/" + events[i].Path
		}
		fsm.collectsCh <- collected
	}
}

func handleCollectedFSEvent(fsm *FSEventMonitor) {
	for collected := range fsm.collectsCh {
		for _, obs := range fsm.collectObservers {
			obs <- collected
		}
	}
}

func NewFSEventMonitor(path string) FSEventMonitor {
	evaluated, err := filepath.EvalSymlinks(path)
	if err != nil {
		glog.Errorf("Error evaluating sym links: %s", err)
		evaluated = path
	}
	return FSEventMonitor{
		startDir:   evaluated,
		collectsCh: make(chan []string),
	}
}

func NewMonitor(startDir string) Monitor {
	fsm := NewFSEventMonitor(startDir)
	return &fsm
}
