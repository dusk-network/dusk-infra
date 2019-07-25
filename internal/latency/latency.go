package latency

import (
	"fmt"
	"io"
	"time"

	"golang.org/x/net/icmp"

	"github.com/sparrc/go-ping"
	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

// Latency carries information about the latency of a Ping to the Voucher Seeder
type Latency struct {
	monitor.Window
	target string //"178.62.193.89"
}

// New creates the Latency monitor.Sampler
func New(t string) *Latency {
	return &Latency{
		Window: make(monitor.Window, 0),
		target: t,
	}
}

// ProbePriviledges checks if the program is allowed to access to raw socket
func (l *Latency) ProbePriviledges() error {
	conn, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func (l *Latency) String() string {
	return "latency"
}

// Monitor writes Param.Value to the websocket
func (l *Latency) Monitor(w io.Writer, m *monitor.Param) error {
	// Pings the voucher seeder
	avgRtt, loss, err := test(l.target)
	if err != nil {
		return err
	}

	if loss > 10 {
		m.Data = map[string]interface{}{"error": fmt.Sprintf("packet loss at %.1f", loss)}
		// TODO: how do we send the notification in this case?
	}

	m.Window = m.Window.Append(float64(avgRtt / time.Millisecond))
	l.Window = l.Add(m.Window)
	if err := j.Write(w, m); err != nil {
		return err
	}
	return nil
}

// InitialState as defined in the StatefulMon interface
func (l *Latency) InitialState(w io.Writer) error {
	if len(l.Window) > 0 {
		m := monitor.NewParam("latency")
		m.Window = l.Window
		if err := j.Write(w, m); err != nil {
			return err
		}
	}
	return nil
}

func test(addr string) (avgRttMs time.Duration, loss float64, err error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return 0, 0, err
	}
	pinger.SetPrivileged(true)
	pinger.Count = 3
	// this is a blocking call
	pinger.Run()
	stats := pinger.Statistics()
	return stats.AvgRtt, stats.PacketLoss, nil
}
