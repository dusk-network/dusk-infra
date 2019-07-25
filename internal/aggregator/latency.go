package aggregator

import (
	"fmt"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeLatency(p monitor.Param) string {
	err, errFound := p.Data["error"]
	if errFound {
		return err.(string)
	}

	w := c.status.latency.Add(p.Window)
	avg := w.CalculateAvg()
	c.lock.Lock()
	c.status.latency = w
	c.status.Latency = avg
	c.lock.Unlock()

	if avg > 150 {
		return fmt.Sprintf("network too slow. Latency more than 150ms (%.0fms)", avg)
	}

	return ""
}
