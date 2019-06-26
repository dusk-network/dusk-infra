package monitor

import (
	"io"
	"time"

	lg "github.com/sirupsen/logrus"
)

var log = lg.WithField("process", "monitor")

type (
	Mon interface {
		Wire(io.WriteCloser)
		Disconnect()
	}

	Param struct {
		Timestamp time.Time              `json:"timestamp"`
		Metric    string                 `json:"metric"`
		Value     string                 `json:"value"`
		Data      map[string]interface{} `json:"data"`
	}

	Supervisor interface {
		Monitor(io.Writer, *Param) error
	}

	TickerMonitor struct {
		Supervisor
		Metric   string
		quitChan chan struct{}
		i        time.Duration
	}
)

func New(s Supervisor, i time.Duration, metric string) *TickerMonitor {
	return &TickerMonitor{
		Supervisor: s,
		i:          i,
		quitChan:   make(chan struct{}, 1),
		Metric:     metric,
	}
}

// Wire a Monitor to a writer. Usually the writer is an outgoing connection. The logic for production and forwarding of the packet is a simple ticker.
func (m *TickerMonitor) Wire(w io.WriteCloser) {
	ticker := time.NewTicker(m.i)
	for {
		select {
		case <-ticker.C:
			if err := m.write(w); err != nil {
				log.WithError(err).Errorln("connection problem")
			}
		case <-m.quitChan:
			log.Infoln("quitting on request of the client")
			_ = w.Close()
			ticker.Stop()
			return
		}
	}
}

func (m *TickerMonitor) Disconnect() {
	m.quitChan <- struct{}{}
}

func (m *TickerMonitor) write(w io.WriteCloser) error {
	p := &Param{
		Metric:    m.Metric,
		Timestamp: time.Now(),
	}
	log.WithField("param", p.Metric).Debugln("packet produced")
	if err := m.Monitor(w, p); err != nil {
		return err
	}
	return nil
}
