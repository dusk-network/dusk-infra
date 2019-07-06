package mux

import (
	"io"
	"math/rand"
	"sync"

	lg "github.com/sirupsen/logrus"

	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
)

var log = lg.WithField("process", "mux")

type Writer struct {
	sync.Mutex
	io.Writer
	conns map[uint32]*j.SynConn
}

func New() *Writer {
	return &Writer{
		conns: make(map[uint32]*j.SynConn),
	}
}

// Write on the websocket
func (mux *Writer) Write(b []byte) (int, error) {
	if mux.Writer != nil {
		return mux.Writer.Write(b)
	}
	log.Debugln("No writer specified yet. Dropping packet")
	return 0, nil
}

// Add a websocket and writes the initial state where appropriate
func (mux *Writer) Add(jw j.JsonReadWriter) uint32 {
	mux.Lock()
	defer mux.Unlock()
	id := rand.Uint32()
	mux.conns[id] = j.New(jw)
	mux.rebuildWriter()
	return id
}

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
