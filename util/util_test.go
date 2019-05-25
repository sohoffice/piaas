package util

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

// The method to setup and tear down the tests of this package.
func TestMain(m *testing.M) {
	debugPtr := flag.Bool("debug", false, "Print debug messsages")
	flag.Parse()

	if *debugPtr {
		log.SetLevel(log.DebugLevel)
	}
	os.Exit(m.Run())
}
