package util

import (
	"testing"
	"time"
)

func TestDebouncer(t *testing.T) {
	ch := make(chan string)
	collected := make([]string, 0)
	go func() {
		msg := <-ch
		switch msg {
		case "":
			collected = append(collected, msg)
		}
	}()
	debouncer := NewDebouncer(100, func() {
		ch <- ""
	})
	// trigger event multiple times
	debouncer.Event()
	debouncer.Event()
	debouncer.Event()

	// wait 150 millis
	<-time.After(time.Millisecond * 110)

	// validate the collection
	if len(collected) != 1 {
		t.Errorf("Collected more than 1 element: %s", collected)
	}
	if collected[0] != "" {
		t.Errorf("Should only collect an empty string: %s", collected)
	}
}
