package util

import (
	"github.com/golang/glog"
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
		default:
			t.Errorf("Unexpected message: %s", msg)
		}
	}()
	debouncer := NewDebouncer(100, func(tag string) {
		glog.Infof("Debounce event received: %s", tag)
		ch <- ""
	})
	// trigger event multiple times
	debouncer.Event("1")
	debouncer.Event("2")
	debouncer.Event("3")

	// wait 150 millis
	<-time.After(time.Millisecond * 210)

	// validate the collection
	if len(collected) != 1 {
		t.Errorf("Should have collected 1 elements: %s", collected)
	}
	if collected[0] != "" {
		t.Errorf("Should only collect an empty string: %s", collected)
	}
}
