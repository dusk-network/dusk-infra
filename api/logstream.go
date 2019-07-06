package logstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	lg "github.com/sirupsen/logrus"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var quitSig = struct{}{}
var log = lg.WithField("process", "logstream")

// LogStreamMonitor is a TCP server listening for json-encoded data from a single (local) connection.
// It exposes an error channel for notifying of problems and listens a quit channel for termination requests
// Packets are json-encoded and are proxied to both the client and an aggregator. Ideally, it would use a stream multiplexer to abstract writers away but given its focus on being a simple channel for node logging, the complication of introducing an abstraction layer is  rather unnecessary
// LogStreamMonitor implements the monitor.Mon interface
type LogStreamMonitor struct {
	srv      net.Listener
	ErrChan  chan error
	quitChan chan struct{}
	dataChan chan monitor.Param
}

// New creates a new LogProxy from a host. The host should be a correct URL (such as unix:///path/to/unix.sock)
func New(h *url.URL) *LogStreamMonitor {
	fmt.Println("Listening on unix socket")
	log.WithField("URL", h.String()).Debugln("starting logstream server")
	srv, err := net.Listen(h.Scheme, h.Path)
	if err != nil {
		log.Panic(err)
	}

	return &LogStreamMonitor{
		srv:      srv,
		ErrChan:  make(chan error),
		quitChan: make(chan struct{}, 1),
		dataChan: make(chan monitor.Param),
	}
}

func (l *LogStreamMonitor) Shutdown() {
	l.quitChan <- struct{}{}
	_ = l.srv.Close()
}

// Pipe a json.Writer to the incoming connection. As a result, the data sent over the line will be encoded in the monitor.Param struct on the `Data` field.
// The Value field will be left empty and the `Metric` set to `log`
func (l *LogStreamMonitor) Wire(w io.Writer) {

	log.Debugln("wiring the logstream monitoring")
	for {
		conn, err := l.srv.Accept()
		if err != nil {
			log.WithError(err).Warnln("error in creating the connection")
			l.ErrChan <- err
			return
		}

		defer conn.Close()

		d := json.NewDecoder(conn)
		go l.receive(d)
		for {
			log.Debug("waiting for packet")
			select {
			case p := <-l.dataChan:
				log.Debugf("got packet: %s\n", p.String())
				param, err := json.Marshal(&p)
				if err != nil {
					log.WithError(err).Warnln("error in package reception")
					l.ErrChan <- err
					return
				}

				if _, err := w.Write(param); err != nil {
					log.WithError(err).Warnln("exiting")
					l.ErrChan <- err
					return
				}
			case <-l.quitChan:
				log.Warnln("quitting")
				return
			}
		}
	}
}

func (l *LogStreamMonitor) receive(d *json.Decoder) {
	for {
		var msg map[string]interface{}
		if err := d.Decode(&msg); err != nil {
			l.ErrChan <- err
			log.WithError(err).Warnln("error in decoding incoming JSON packet")
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

func (l *LogStreamMonitor) Disconnect() {
	l.quitChan <- quitSig
}
