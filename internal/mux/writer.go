package mux

import (
	"bytes"
	"io"
	"math/rand"
	"sync"

	lg "github.com/sirupsen/logrus"

	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
)

var log = lg.WithField("process", "mux")

// Writer is a threadsafe multiwriter to the websocket connections
type Writer struct {
	sync.Mutex
	io.Writer
	conns map[uint32]*j.SynConn
}

// New creates a Writer
func New() *Writer {
	return &Writer{
		conns: make(map[uint32]*j.SynConn),
	}
}

// Write on the websocket
func (mux *Writer) Write(b []byte) (int, error) {
	writer := mux.Writer
	if writer == nil {
		// faking a Writer so samplers will start to build their windows
		writer = new(bytes.Buffer)
		log.Debugln("No real writer specified yet. Buffering windows")
	}

	return writer.Write(b)
}

// Add a websocket and writes the initial state where appropriate
func (mux *Writer) Add(jw j.ReadWriter) uint32 {
	mux.Lock()
	defer mux.Unlock()
	id := rand.Uint32()
	mux.conns[id] = j.New(jw)
	mux.rebuildWriter()
	return id
}

// Remove a websocket from the Writer
func (mux *Writer) Remove(id uint32) {
	mux.Lock()
	defer mux.Unlock()
	delete(mux.conns, id)
	mux.rebuildWriter()
}

func (mux *Writer) rebuildWriter() {
	var curConn []io.Writer
	for _, conn := range mux.conns {
		curConn = append(curConn, conn)
	}

	mux.Writer = io.MultiWriter(curConn...)
}
