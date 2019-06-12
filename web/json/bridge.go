package json

import (
	"io"
	"sync"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

// convenience function that mimics the websocket.Conn.WriteJSON
type Writer interface {
	WriteJSON(v interface{}) error
}

type Bridge struct {
	lock sync.Mutex
	jw   Writer
	m    *monitor.Monitor
}

func New(c Writer, s *monitor.Monitor) *Bridge {
	return &Bridge{
		lock: sync.Mutex{},
		jw:   c,
		m:    s,
	}
}

// Connect is used to connect the monitor outcome to the websocket
func (jb *Bridge) Connect() error {
	r, w := io.Pipe()

	// wiring the Supervisor to a pipe to decouple the two
	go jb.m.Wire(w)

	for {
		b := make([]byte, 0)
		if _, err := r.Read(b); err != nil {
			// EOF means disconnect
			if err == io.EOF {
				return nil
			}
			return err
		}

		// forwarding the message to the websocket connection and locking the connection at the same time
		if err := jb.WriteJSON(string(b)); err != nil {
			jb.lock.Unlock()
			return err
		}
	}
}

func (jb *Bridge) Disconnect() {
	jb.m.QuitChan <- struct{}{}
}

// WriteJSON writes the JSON encoding of v as a message.
func (jb *Bridge) WriteJSON(v string) error {
	jb.lock.Lock()
	defer jb.lock.Unlock()
	return jb.jw.WriteJSON(v)
}
