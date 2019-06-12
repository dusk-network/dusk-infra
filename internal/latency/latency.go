package latency

import (
	"encoding/json"
	"fmt"
	"io"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/ping"
)

type Latency struct {
	target string //"178.62.193.89"
}

func New(t string) *Latency {
	return &Latency{target: t}
}

func (l *Latency) Monitor(w io.Writer, m *monitor.Param) error {

	// Pings the voucher seeder
	delay, err := test(l.target)
	if err != nil {
		return err
	}
	m.Value = fmt.Sprintf("%f", delay)
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func test(addr string) (float64, error) {
	_, dur, err := ping.Ping(addr)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return (float64(dur) / 1000000), nil
}
