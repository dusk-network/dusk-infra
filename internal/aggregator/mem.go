package aggregator

import (
	"fmt"
	"strconv"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeMem(p *monitor.Param) string {
	mem, err := strconv.ParseFloat(p.Value, 32)
	if err != nil {
		log.WithError(err).Warnln("error in parsing the mem value")
	}

	c.status.Mem = float32(mem)
	if mem > 80 {
		return fmt.Sprintf("high memory usage (%s%%)", p.Value)
	}
	return ""
}
