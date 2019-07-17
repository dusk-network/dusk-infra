package aggregator

import (
	"fmt"
	"strconv"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeCPU(p *monitor.Param) string {
	cpu, err := strconv.ParseFloat(p.Value, 64)
	if err != nil {
		log.WithError(err).Warnln("error in parsing the cpu value")
		return ""
	}

	c.lock.Lock()
	w := c.status.cpu.Append(cpu)
	avg := w.CalculateAvg()
	c.status.cpu = w
	c.status.CPU = avg
	c.lock.Unlock()
	if avg > 50 {
		return fmt.Sprintf("high CPU load (%.2f%%)", avg)
	}
	return ""
}
