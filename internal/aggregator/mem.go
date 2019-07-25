package aggregator

import (
	"fmt"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeMem(p monitor.Param) string {
	w := c.status.mem.Add(p.Window)
	avg := w.CalculateAvg()
	c.lock.Lock()
	c.status.mem = w
	c.status.Mem = avg
	c.lock.Unlock()
	if avg > 80 {
		return fmt.Sprintf("high memory usage (%.2f%%)", avg)
	}
	return ""
}
