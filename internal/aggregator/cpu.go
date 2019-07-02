package aggregator

import (
	"fmt"
	"strconv"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeCpu(p *monitor.Param) string {
	cpu, err := strconv.ParseFloat(p.Value, 32)
	if err != nil {
		log.WithError(err).Warnln("error in parsing the cpu value")
	}

	if cpu > 55 {
		return fmt.Sprintf("high CPU load (%s%%)", p.Value)
	}
	return ""
}
