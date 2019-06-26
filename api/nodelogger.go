package nodelogger

import (
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"time"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
	mon "gitlab.dusk.network/dusk-core/node-monitor/web/json"
)

var quitSig = struct{}{}

// LogProxy is a TCP server listening for json-encoded data from a single (local) connection.
// It exposes an error channel for notifying of problems and listens a quit channel for termination requests
// Packets are json-encoded and are proxied to both the client and an aggregator. Ideally, it would use a stream multiplexer to abstract writers away but given its focus on being a simple channel for node logging, the complication of introducing an abstraction layer is  rather unnecessary
// LogProxy implements the monitor.Mon interface
type LogProxy struct {
	host *url.URL
	// aggregator *url.URL
	ErrChan  chan error
	quitChan chan struct{}
	dataChan chan monitor.Param
}

// New creates a new LogProxy from a host. The host should be a correct URL (such as unix:///path/to/unix.sock)
func New(h *url.URL) *LogProxy {
	return &LogProxy{
		host:     h,
		ErrChan:  make(chan error),
		quitChan: make(chan struct{}),
		dataChan: make(chan monitor.Param),
	}
}

// Pipe a json.Writer to the incoming connection. As a result, the data sent over the line will be encoded in the monitor.Param struct on the `Data` field.
// The Value field will be left empty and the `Metric` set to `log`
func (l *LogProxy) Pipe(w mon.Writer) {
	srv, err := net.Listen(l.host.Scheme, l.host.Path)
	if err != nil {
		l.ErrChan <- err
		return
	}

	conn, err := srv.Accept()
	if err != nil {
		l.ErrChan <- err
		return
	}

	defer conn.Close()

	if conn == nil {
		l.ErrChan <- errors.New("connection is nil")
		return
	}
	// io.TeeReader()

	d := json.NewDecoder(conn)
	go l.receive(d)
	for {
		select {
		case p := <-l.dataChan:
			strJson, err := json.Marshal(&p)
			if err != nil {
				l.ErrChan <- err
				return
			}

			if err := w.WriteJSON(string(strJson)); err != nil {
				l.ErrChan <- err
				return
			}
		case <-l.quitChan:
			return
		}
	}
}

func (l *LogProxy) receive(d *json.Decoder) {
	for {
		var msg map[string]interface{}
		if err := d.Decode(&msg); err != nil {
			l.ErrChan <- err
			return
		}
		// forwarding the message to the websocket connection and locking the connection at the same time
		l.dataChan <- monitor.Param{
			Metric:    "log",
			Timestamp: time.Now(),
			Data:      msg,
		}
	}
}

func (l *LogProxy) Disconnect() {
	l.quitChan <- quitSig
}
