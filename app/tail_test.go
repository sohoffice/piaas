package app

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/util"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

// Prepare a file with 10 lines
func prepare() (string, string, string) {
	tempDir := os.TempDir()
	var buf bytes.Buffer

	for i := 0; i < 10; i++ {
		buf.WriteString(fmt.Sprintf("line %d\n", i+1))
	}
	s := buf.String()

	filename := path.Join(tempDir, "a.txt")
	err := ioutil.WriteFile(filename, []byte(s), 0644)
	util.CheckError("write test file", err)

	return tempDir, filename, buf.String()
}

func TestTail(t *testing.T) {
	tempDir, filename, content := prepare()
	defer os.RemoveAll(tempDir)
	var buffer bytes.Buffer

	tail := NewTail(filename, &buffer, 10)
	tail.Start()
	tail.Read(3)
	<-time.After(time.Millisecond * 200)
	log.Infof("Last 3 lines: %s", buffer.String())

	if buffer.String() != "line 8\nline 9\nline 10\n" {
		t.Errorf("Should read the last 3 lines, but read: %s", buffer.String())
	}
	if int(tail.pos) != len(content) {
		t.Errorf("Should advanced position to %d, but is %d", len(content), tail.pos)
	}

	buffer.Reset()
	tail.pos = 0
	tail.Read(1)
	<-time.After(time.Millisecond * 200)
	log.Infof("Last 1 line: %s", buffer.String())
	if buffer.String() != "line 10\n" {
		t.Errorf("Should read the last 1 line, but read: %s", buffer.String())
	}

	buffer.Reset()
	tail.pos = 49
	tail.Read(-1)
	<-time.After(time.Millisecond * 200)
	log.Infof("Skip the first 49 bytes: %s", buffer.String())
	//<- time.After(time.Minute * 5)
	if buffer.String() != "line 8\nline 9\nline 10\n" {
		t.Errorf("Should skip the first 49 bytes and read the last 3 lines, but read: %s", buffer.String())
	}

	buffer.Reset()
	tail.pos = 0
	tail.Read(20)
	<-time.After(time.Millisecond * 400)
	// <- time.After(time.Minute * 5)
	if buffer.String() != content {
		t.Errorf("Should read the entire file, but read:\n%s", buffer.String())
	}
}
