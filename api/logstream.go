package logstream

import (
	"encoding/json"
	"io"
	"net"
	"net/url"
	"os"
	"time"

	lg "github.com/sirupsen/logrus"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var quitSig = struct{}{}
var log = lg.WithField("process", "logstream")

// status saves the last Param per code to send it upon a new incoming
// connection
type status monitor.Param

// Tolerance represents the maximum amount of time to retain data
var Tolerance = time.Minute * 15

// Monitor is a TCP server listening for json-encoded data from a single (local) connection.
// It exposes an error channel for notifying of problems and listens a quit channel for termination requests
// Packets are json-encoded and are proxied to both the client and an aggregator. Ideally, it would use a stream multiplexer to abstract writers away but given its focus on being a simple channel for node logging, the complication of introducing an abstraction layer is  rather unnecessary
// Monitor implements the monitor.Mon interface
type Monitor struct {
	state    status
	srv      net.Listener
	ErrChan  chan error
	quitChan chan struct{}
	dataChan chan monitor.Param
}

func (s status) merge(p monitor.Param) status {
	code, found := p.Data["code"]
	if !found {
		return s
	}

	switch code {
	case "goroutine":
		s = s.mergeValue(p.Data, "nr", "threads")
	case "round":
		s = s.mergeValue(p.Data, "blockTime", "blockTimes")
		s.Data["round"] = p.Data["round"]
		s.Data["blockHash"] = p.Data["blockHash"]
	default:
		log.WithField("code", code).Warnln("unrecognized code")
	}
	return s
}

func (s status) mergeValue(data map[string]interface{}, k, target string) status {
	v, found := data[k]
	if !found {
		return s
	}

	val, ok := v.(float64)
	if !ok {
		return s
	}

	vals, there := s.Data[target]
	if !there {
		vals = monitor.NewWindow()
	}

	valWindow := vals.(monitor.Window)
	valWindow = valWindow.Append(val)
	s.Data[target] = valWindow
	return s
}

// New creates a new LogProxy from a host. The host should be a correct URL (such as unix:///path/to/unix.sock)
func New(h *url.URL) *Monitor {
	log.WithField("URL", h.String()).Debugln("starting logstream server")
	srv, err := net.Listen(h.Scheme, h.Path)
	if err != nil {
		_ = os.Remove(h.Path)
		srv, err = net.Listen(h.Scheme, h.Path)
		if err != nil {
			panic(err)
		}
	}

	p := status(*monitor.NewParam("status"))

	return &Monitor{
		state:    p,
		srv:      srv,
		ErrChan:  make(chan error),
		quitChan: make(chan struct{}, 1),
		dataChan: make(chan monitor.Param),
	}
}

// Shutdown the server and sends a message to the quit channel
func (l *Monitor) Shutdown() {
	l.quitChan <- struct{}{}
	_ = l.srv.Close()
}

// Wire a json.Writer to the incoming connection. As a result, the data sent over the line will be encoded in the monitor.Param struct on the `Data` field.
// The Value and Window fields will be left empty and the `Metric` set to `log`
func (l *Monitor) Wire(multiwrt io.Writer) {

	log.Debugln("wiring the logstream monitoring")
	for {
		uxconn, err := l.srv.Accept()
		log.Debugln("New incoming client")
		if err != nil {
			log.WithError(err).Warnln("error in creating the connection")
			l.ErrChan <- err
			return
		}

		defer uxconn.Close()

		go func(c net.Conn) {
			d := json.NewDecoder(c)
			go l.receive(d)
			for {
				log.Debug("waiting for packet")
				select {
				case p := <-l.dataChan:
					log.Debugf("got packet: %s\n", p.String())
					if err := l.forward(multiwrt, p); err != nil {
						return
					}
					l.state = l.state.merge(p)
				case <-l.quitChan:
					log.Warnln("quitting")
					return
				}
			}
		}(uxconn)
	}
}

// String is called by the server when connecting the Websocket to log the name
// of the current Sampler
func (l *Monitor) String() string {
	return "logstream"
}

// InitialState as specified by the monitor.StatefulMon interface
func (l *Monitor) InitialState(multiwrt io.Writer) error {
	p := monitor.Param(l.state)
	return l.forward(multiwrt, p)
}

func (l *Monitor) forward(multiwrt io.Writer, p monitor.Param) error {
	log.WithField("param", p).Traceln("sending param")
	param, err := json.Marshal(p)
	if err != nil {
		log.WithError(err).Warnln("error in package reception")
		l.ErrChan <- err
		return err
	}

	if _, err := multiwrt.Write(param); err != nil {
		log.WithError(err).Warnln("exiting")
		l.ErrChan <- err
		return err
	}
	log.Traceln("sent param")
	return nil
}

func (l *Monitor) receive(d *json.Decoder) {
	for {
		var msg map[string]interface{}
		if err := d.Decode(&msg); err != nil {
			l.ErrChan <- err
			log.WithError(err).Warnln("error in decoding incoming JSON packet")
			return
		}
		if len(msg) > 0 {

			// forwarding the message to the websocket connection and locking the connection at the same time
			l.dataChan <- monitor.Param{
				Metric:    "log",
				Timestamp: time.Now(),
				Data:      msg,
			}
		}
	}
}

// Disconnect the monitoring
func (l *Monitor) Disconnect() {
	l.quitChan <- quitSig
}
