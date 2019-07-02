package aggregator

import (
	"fmt"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeLog(p *monitor.Param) string {
	var payload string

	switch p.Data["code"] {
	case "round":
		round := p.Data["round"]
		hash := p.Data["blockHash"]
		time := p.Data["blockTime"]

		payload = fmt.Sprintf("new block validated: round %s, hash %s, block time %s", round, hash, time)
	case "warn":
		level := p.Data["level"]
		msg := p.Data["msg"]
		payload = fmt.Sprintf("[%s] %s", level, msg)
	}

	return payload
}