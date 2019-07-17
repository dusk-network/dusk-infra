package aggregator

import (
	"fmt"
	"strconv"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeLatency(p *monitor.Param) string {
	err, errFound := p.Data["error"]
	if errFound {
		return err.(string)
	}

	v, er := strconv.ParseFloat(p.Value, 64)
	if er != nil {
		log.WithError(er).Warnln("latency processing reported error")
		return ""
	}
	c.lock.Lock()
	latency := c.status.latency.Append(v)
	avg := latency.CalculateAvg()
	c.status.latency = latency
	c.status.Latency = avg
	c.lock.Unlock()
	if avg > 150 {
		return fmt.Sprintf("network too slow. Latency more than 150ms (%.0fms)", avg)
	}

	return ""
}
