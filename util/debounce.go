package util

import "time"

// Signal a potential start of the debounce monitor
type DebounceEvent bool

// handle any message being posted to the debounce channel and start a new debounce period if no current period is running.
type Debouncer struct {
	// The length of debounce period in milliseconds.
	millis int64
	// the function to call when debounce period completes
	notify func()
	// whether a current debounce period is running
	running bool
	// the channel to receive start request of debounce period.
	// Any Debounceevent can be sent to trigger, usually a true.
	ch chan DebounceEvent
}

// Send an event to request to start or join a debounce period
func (dbPtr *Debouncer) Event() {
	(*dbPtr).ch <- true
}

func handleDebounceEvent(dbPtr *Debouncer) {
	db := *dbPtr
	for {
		select {
		case _ = <-db.ch:
			if !db.running {
				go func() {
					db.running = true
					// wait for debounce millis
					<-time.After(time.Millisecond * time.Duration(db.millis))

					db.notify()
					db.running = false
				}()
			}
		}
	}
}

// Create a new debouncer
func NewDebouncer(millis int64, notify func()) Debouncer {
	db := Debouncer{
		millis:  millis,
		notify:  notify,
		running: false,
		ch:      make(chan DebounceEvent),
	}
	go handleDebounceEvent(&db)

	return db
}
