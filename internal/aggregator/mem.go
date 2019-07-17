package aggregator

import (
	"fmt"
	"strconv"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeMem(p *monitor.Param) string {
	mem, err := strconv.ParseFloat(p.Value, 64)
	if err != nil {
		log.WithError(err).Warnln("error in parsing the mem value")
		return ""
	}

	c.lock.Lock()
	m := c.status.mem.Append(mem)
	avg := m.CalculateAvg()
	c.status.Mem = avg
	c.status.mem = m
	c.lock.Unlock()
	if avg > 80 {
		return fmt.Sprintf("high memory usage (%.2f%%)", avg)
	}
	return ""
}
