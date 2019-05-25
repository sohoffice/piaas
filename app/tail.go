package app

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/util"
	"io"
	"os"
)

type Tail struct {
	reader     io.ReadSeeker
	readCloser io.Closer
	writer     io.Writer
	// file content will flow through this channel
	dataCh chan []byte
	// read request will be added to this channel
	readCh  chan int32
	pos     int64
	bufsize int64
	buffer  []byte
}

func NewTail(filename string, writer io.Writer, bufsize int64) *Tail {
	// The channel will
	dataCh := make(chan []byte, 10)
	input, err := os.Open(filename)
	util.CheckFatal("open file for tailing", err)

	go func() {
		for {
			data := <-dataCh
			_, err := writer.Write(data)
			util.CheckError("Write data", err)
		}
	}()

	return &Tail{
		reader:     io.ReadSeeker(input),
		readCloser: input,
		writer:     writer,
		dataCh:     dataCh,
		readCh:     make(chan int32, 10),
		bufsize:    bufsize,
		buffer:     make([]byte, bufsize),
	}
}

// Actually start taking read request.
// The caller is responsible for calling Read() method.
func (t *Tail) Start() {
	go func() {
		for {
			n := <-t.readCh
			if n != 0 { // if n == 0, we don't need to do anything.

				// move to the end of file to mark the end position
				endPos, err := t.reader.Seek(0, 2)
				util.CheckError("seek file end", err)

				var data []byte
				if n < 0 {
					data = t.readAllBytes(endPos)
				} else {
					data = t.readLines(n, endPos)
				}
				t.pos = endPos
				t.dataCh <- data
			}
		}
	}()
}

func (t *Tail) Read(lines int32) {
	t.readCh <- lines
}

func (t *Tail) Close() {
	if t.readCh != nil {
		close(t.readCh)
	}
	if t.dataCh != nil {
		close(t.dataCh)
	}
	if t.readCloser != nil {
		t.readCloser.Close()
	}
}

var newline = []byte("\n")

// Read all data from current position to endPos
func (t *Tail) readAllBytes(endPos int64) []byte {
	var err error
	if endPos <= t.pos { // nothing to read.
		return t.buffer[0:0]
	}
	lengthToRead := endPos - t.pos
	_, err = t.reader.Seek(t.pos, 0)
	util.CheckError("seek file", err)
	// read into the buffer
	readlen, err := t.reader.Read(t.buffer)
	util.CheckError("read file", err)

	// We have read too much.
	if readlen > int(lengthToRead) {
		readlen = int(lengthToRead)
	}

	// we have read enough
	if readlen >= int(lengthToRead) {
		return t.buffer[0:readlen]
	} else {
		var data bytes.Buffer
		data.Write(t.buffer[0:readlen])
		t.pos = t.pos + int64(readlen)
		return append(data.Bytes(), t.readAllBytes(endPos)...)
	}
}

// Read data from the endPos, until n lines were read.
func (t *Tail) readLines(lines int32, endPos int64) []byte {
	buffer := t.buffer
	var err error

	if lines <= 0 || endPos <= 0 {
		return buffer[0:0]
	}

	startPos := endPos - t.bufsize
	if startPos < 0 {
		startPos = 0
	}
	maxReadLen := int(endPos - startPos)

	// rewind a little bit to start reading
	_, err = t.reader.Seek(startPos, 0)
	util.CheckError("seek file", err)
	// read into the buffer
	readlen, err := t.reader.Read(buffer)
	util.CheckError("read file", err)
	if readlen > maxReadLen {
		readlen = maxReadLen
	}
	if readlen > 0 && buffer[readlen-1] == newline[0] {
		// we have an ending newline, increase the required lines count by 1
		lines = lines + 1
	}

	var linesRead int32 = 0
	var newlinePos = readlen + 1
	for linesRead < lines && newlinePos > 0 {
		newlinePos = bytes.LastIndex(buffer[0:newlinePos-1], newline)
		if newlinePos != -1 {
			linesRead = linesRead + 1
			log.Debugf("Find newline: %d, linesRead: %d, lines: %d", newlinePos, linesRead, lines)
		}
	}

	if lines > linesRead { // Not enough lines, so we move to the beginning of the buffer.
		var data bytes.Buffer
		data.Write(buffer[0:readlen])
		log.Debugf("Not enough line, %d more expected: %s", lines-linesRead, data.String())
		return append(t.readLines(lines-linesRead, endPos-int64(readlen)), data.Bytes()...)
	} else if lines == linesRead { // We have exactly what we want. return
		return buffer[newlinePos+1 : readlen]
	} else {
		log.Errorf("Tail error, asking for %d lines, but %d lines were read.", lines, linesRead)
		return buffer[0:0]
	}
}
