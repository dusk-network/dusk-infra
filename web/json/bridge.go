package json

import (
	"bytes"
	"io"

	lg "github.com/sirupsen/logrus"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var log = lg.WithField("process", "json bridge")

// Writer is a convenience interface that mimics the websocket.Conn.WriteJSON.
// It is used to threat websocket as an implementation detail and facilitate testing
type Writer interface {
	WriteJSON(string) error
}

type Bridge struct {
	jw Writer
	m  monitor.Mon
}

func New(c Writer, s monitor.Mon) *Bridge {
	return &Bridge{
		jw: c,
		m:  s,
	}
}

// Connect is used to connect the monitor outcome to the websocket
func (jb *Bridge) Connect() {
	r, w := io.Pipe()

	log.Debugln("bridge connected")
	// wiring the Supervisor to a pipe to decouple the two
	go jb.m.Wire(w)

	for {
		b := make([]byte, 512)
		if _, err := r.Read(b); err != nil {
			// EOF means disconnect
			if err == io.EOF {
				return
			}

			log.WithError(err).Errorln("error in reading from the monitor. Disposing connection")
			return
		}
		log.Debugln("before WriteJSON")
		b = bytes.Trim(b, "\x00")

		// forwarding the message to the websocket connection and locking the connection at the same time
		if err := jb.jw.WriteJSON(string(b)); err != nil {
			log.WithError(err).Errorln("error in writing json encoded data to the writer. Disposing connection")
			return
		}
	}
}
