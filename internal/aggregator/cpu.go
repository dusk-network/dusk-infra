package aggregator

import (
	"fmt"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeCPU(p monitor.Param) string {
	w := c.status.cpu.Add(p.Window)
	avg := w.CalculateAvg()
	c.lock.Lock()
	c.status.cpu = w
	c.status.CPU = avg
	c.lock.Unlock()
	if avg > 50 {
		return fmt.Sprintf("high CPU load (%.2f%%)", avg)
	}
	return ""
}
