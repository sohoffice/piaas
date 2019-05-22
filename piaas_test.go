package piaas

import (
	"flag"
	"os"
	"testing"
)

// The method to setup and tear down the tests of this package.
func TestMain(m *testing.M) {
	flag.Parse()
	// flag.Lookup("logtostderr").Value.Set("true")
	os.Exit(m.Run())
}
