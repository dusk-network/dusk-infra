package latency

import (
	"encoding/json"
	"io"
	"time"

	"golang.org/x/net/icmp"

	"github.com/sparrc/go-ping"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

type Latency struct {
	target string //"178.62.193.89"
}

func New(t string) monitor.Supervisor {
	return &Latency{target: t}
}

func (l *Latency) ProbePriviledges() error {
	conn, err := icmp.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func (l *Latency) Monitor(w io.Writer, m *monitor.Param) error {
	// Pings the voucher seeder
	avgRtt, loss, err := test(l.target)
	if err != nil {
		return err
	}

	if loss > 10 {
		// TODO: how do we send the notification in this case?
	}

	m.Value = avgRtt.String()
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

// func test(addr string) (float64, error) {
// 	_, dur, err := ping.Ping(addr)
// 	if err != nil {
// 		fmt.Println(err)
// 		return 0, err
// 	}
// 	return (float64(dur) / 1000000), nil
// }

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
