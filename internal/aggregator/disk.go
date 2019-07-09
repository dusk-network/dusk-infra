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
	}

	c.lock.Lock()
	c.status.Disk = value
	c.lock.Unlock()
	if value > 90 {
		return fmt.Sprintf("little or no Disk space left (%s%%)", p.Value)
	}
	return ""
}
