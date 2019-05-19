package piaas

import (
	"flag"
	"github.com/golang/glog"
	"os"
	"testing"
)

// The method to setup and tear down the tests of this package.
func TestMain(m *testing.M) {
	defer glog.Flush()
	flag.Parse()
	// flag.Lookup("logtostderr").Value.Set("true")
	os.Exit(m.Run())
}
