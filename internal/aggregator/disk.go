package aggregator

import (
	"fmt"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeDisk(p monitor.Param) string {
	w := c.status.disk.Add(p.Window)
	avg := w.CalculateAvg()
	c.lock.Lock()
	c.status.disk = w
	c.status.Disk = avg
	c.lock.Unlock()
	if avg > 90 {
		return fmt.Sprintf("little or no Disk space left (%.2f%%)", avg)
	}
	return ""
}
