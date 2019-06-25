package nodelogger

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/url"
	"time"

	lg "github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
	mon "gitlab.dusk.network/dusk-core/node-monitor/web/json"
)

var log = lg.WithField("process", "log_proxy")

type LogProxy struct {
	host *url.URL
	// aggregator *url.URL
	QuitChan chan error
}

func New(host string) *LogProxy {
	h, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	return &LogProxy{
		host:     h,
		QuitChan: make(chan error),
	}
}

func (l *LogProxy) Pipe(w mon.Writer, res chan bool) {
	srv, err := net.Listen(l.host.Scheme, l.host.Path)
	if err != nil {
		l.QuitChan <- err
		return
	}

	res <- true
	conn, err := srv.Accept()
	if err != nil {
		l.QuitChan <- err
		return
	}

	defer conn.Close()

	if conn == nil {
		l.QuitChan <- errors.New("connection is nil")
		return
	}
	// io.TeeReader()

	d := json.NewDecoder(conn)
	for {
		var msg map[string]interface{}
		if err := d.Decode(&msg); err == io.EOF {
			return
		} else if err != nil {
			l.QuitChan <- err
			return
		}

		// forwarding the message to the websocket connection and locking the connection at the same time
		p := &monitor.Param{
			Metric:    "log",
			Timestamp: time.Now(),
			Data:      msg,
		}

		strJson, err := json.Marshal(p)
		if err != nil {
			l.QuitChan <- err
			return
		}

		if err := w.WriteJSON(string(strJson)); err != nil {
			log.WithError(err).Errorln("error in writing json encoded data to the writer. Disposing connection")
			return
		}
		res <- true
	}
}
