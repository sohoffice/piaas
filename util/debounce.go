package util

import (
	log "github.com/sirupsen/logrus"
	"time"
)

// handle any message being posted to the debounce channel and start a new debounce period if no current period is running.
type Debouncer struct {
	// The length of debounce period in milliseconds.
	millis uint64
	// the function to call when debounce period completes
	notify func(string)
	// whether a current debounce period is running
	running bool
	// the channel to receive start request of debounce period.
	// Any Debounceevent can be sent to trigger, usually a true.
	ch chan string
}

// Send an event to request to start or join a debounce period
func (dbPtr *Debouncer) Event(tag string) {
	(*dbPtr).ch <- tag
}

func handleDebounceEvent(dbPtr *Debouncer) {
	db := *dbPtr
	for {
		tag := <-db.ch
		// log.Debugf("Debouncer triggered: %s, %t", tag, db.running)
		if !db.running {
			log.Debugf("Start debounce watcher: %s", tag)
			db.running = true
			go func() {
				// wait for debounce millis
				<-time.After(time.Millisecond * time.Duration(db.millis))

				db.notify(tag)
				db.running = false
			}()
		}
	}
}

// Create a new debouncer
func NewDebouncer(millis uint64, notify func(string)) Debouncer {
	db := Debouncer{
		millis:  millis,
		notify:  notify,
		running: false,
		ch:      make(chan string),
	}
	go handleDebounceEvent(&db)

	return db
}
