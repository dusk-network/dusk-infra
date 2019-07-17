package aggregator

import (
	"fmt"
	"strconv"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeDisk(p *monitor.Param) string {
	value, err := strconv.ParseFloat(p.Value, 64)
	if err != nil {
		log.WithError(err).Warnln("error in parsing the disk value")
		return ""
	}

	c.lock.Lock()
	disk := c.status.disk.Append(value)
	avg := disk.CalculateAvg()
	c.status.disk = disk
	c.status.Disk = avg
	c.lock.Unlock()
	if avg > 90 {
		return fmt.Sprintf("little or no Disk space left (%.2f%%)", avg)
	}
	return ""
}
