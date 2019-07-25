package json

import (
	j "encoding/json"
	"io"
	"sync"

	lg "github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var log = lg.WithField("process", "SynConn")

// MessageReader is the Reader interface of WebSocket
type MessageReader interface {
	ReadMessage() (int, []byte, error)
}

// Writer extrapolates the json Writer interface of WebSocket
type Writer interface {
	WriteJSON(v interface{}) error
	Close() error
}

// ReadWriter extrapolates the json ReaderWriter interface of WebSocket
type ReadWriter interface {
	MessageReader
	Writer
}

// SynConn is the shared (Websocket) connection among all processes that
// independently write on it
type SynConn struct {
	sync.RWMutex
	ReadWriter
}

// New creates a SynConn from a JSONReadWriter (Websocket)
func New(w ReadWriter) *SynConn {
	return &SynConn{
		ReadWriter: w,
	}
}

// Write is a convenience wrapper to allow locking of the connection when a
// process is writing on it
func (s *SynConn) Write(b []byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	p := &monitor.Param{}

	if err := j.Unmarshal(b, p); err != nil {
		return 0, err
	}

	log.WithField("param", p).Traceln("Unmarshalled JSON into Parameter")
	if err := s.WriteJSON(p); err != nil {
		return 0, err
	}

	log.WithField("param", p).Traceln("Written to SynConn")
	return len(b), nil
}

// ReadMessage locks the websocket when reading
func (s *SynConn) ReadMessage() (int, []byte, error) {
	log.Debugln("reading incoming messages")
	s.RLock()
	defer s.RUnlock()
	return s.ReadWriter.ReadMessage()
}

// Write is a convenience function to write a json-encoded Param
func Write(w io.Writer, m *monitor.Param) error {
	b, err := j.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
