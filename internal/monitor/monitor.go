package monitor

import (
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	Param struct {
		Timestamp time.Time `json:"timestamp"`
		Metric    string    `json:"metric"`
		Value     string    `json:"value"`
	}

	Supervisor interface {
		Monitor(io.Writer, *Param) error
	}

	Monitor struct {
		Supervisor
		Metric   string
		QuitChan chan struct{}
		i        time.Duration
	}
)

func New(s Supervisor, i time.Duration) *Monitor {
	return &Monitor{
		Supervisor: s,
		i:          i,
		QuitChan:   make(chan struct{}),
	}
}

func (m *Monitor) Wire(w io.WriteCloser) {
	for {
		select {
		case <-time.After(m.i):
			p := &Param{
				Metric:    m.Metric,
				Timestamp: time.Now(),
			}
			if err := m.Monitor(w, p); err != nil {
				log.WithError(err).Errorln("connection problem")
				continue
			}
		case <-m.QuitChan:
			log.Infoln("quitting on request of the client")
			_ = w.Close()
			return
		}
	}
}
