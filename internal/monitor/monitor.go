package monitor

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	lg "github.com/sirupsen/logrus"
)

var log = lg.WithField("process", "monitor")

type (
	// StatefulMon is a monitoring process that carries an initial state
	// to be communicated to new inbound connections
	StatefulMon interface {
		Mon
		StatefulSampler
	}

	// StatefulSampler extends the Sampler interface with a state
	StatefulSampler interface {
		InitialState(io.Writer) error
	}

	// Mon is the monitoring process
	Mon interface {
		Wire(io.Writer)
		Shutdown()
	}

	// Param is the json encodable structure communicated to the monitoring clients
	Param struct {
		Timestamp time.Time              `json:"timestamp"`
		Metric    string                 `json:"metric"`
		Window    Window                 `json:"slice,omitempty"`
		Data      map[string]interface{} `json:"data,omitempty"`
		Value     string                 `json:"text,omitempty"` //Deprecated, it should ideally be a []string
	}

	// Sampler wraps the various monitoring processes
	Sampler interface {
		Monitor(io.Writer, *Param) error
	}

	// TickerMonitor pushes data collected from the monitoring with a given frequency
	TickerMonitor struct {
		Sampler
		Metric   string
		quitChan chan struct{}
		i        time.Duration
	}
)

// NewParam builds a new Param
func NewParam(metric string) *Param {
	return &Param{
		Timestamp: time.Now(),
		Metric:    metric,
		Data:      make(map[string]interface{}),
		Window:    NewWindow(),
		Value:     "",
	}
}

func (p *Param) String() string {
	if b, err := json.Marshal(p); err == nil {
		return string(b)
	}

	return ""
}

// Add a value to the Value field
func (p *Param) Add(v float64) {
	p.Window.Append(v)
}

// New creates a TickerMonitor
func New(s Sampler, i time.Duration, metric string) *TickerMonitor {
	return &TickerMonitor{
		Sampler:  s,
		i:        i,
		quitChan: make(chan struct{}, 1),
		Metric:   metric,
	}
}

// Wire a Monitor to a writer. Usually the writer is an outgoing connection. The logic for production and forwarding of the packet is a simple ticker.
func (m *TickerMonitor) Wire(w io.Writer) {
	ticker := time.NewTicker(m.i)
	for {
		select {
		case <-ticker.C:
			if err := m.write(w); err != nil {
				log.WithError(err).Errorln("connection problem")
			}
		case <-m.quitChan:
			log.WithField("metric", m.Metric).Infoln("quitting on request of the client")
			ticker.Stop()
			return
		}
	}
}

// Shutdown sends a signal to the internal quit channel
func (m *TickerMonitor) Shutdown() {
	m.quitChan <- struct{}{}
}

func (m *TickerMonitor) write(w io.Writer) error {
	p := NewParam(m.Metric)

	log.WithField("param", p.Metric).Debugln("packet produced")
	if err := m.Monitor(w, p); err != nil {
		return err
	}
	return nil
}

func (m *TickerMonitor) String() string {
	return fmt.Sprintf("%s", m.Sampler)
}

// InitialState as specified by the StatefulMon interface
func (m *TickerMonitor) InitialState(w io.Writer) error {
	initializer, ok := m.Sampler.(StatefulSampler)
	if ok {
		return initializer.InitialState(w)
	}
	return nil
}
