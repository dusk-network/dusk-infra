package aggregator

import (
	"fmt"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func (c *Client) serializeLog(p *monitor.Param) (string, string) {
	var payload string
	code := p.Data["code"]

	switch code {
	case "round":
		round := p.Data["round"]
		hash := p.Data["blockHash"]
		time := p.Data["blockTime"]
		c.lock.Lock()
		c.status.BlockHash = hash.(string)
		c.status.BlockTime = time.(float64)
		c.status.Round = uint64(round.(float64))
		c.lock.Unlock()

		payload = fmt.Sprintf("new block validated: round %d, hash %s, block time %.2fms", c.status.Round, hash, time)
	case "warn":
		level := p.Data["level"]
		msg := p.Data["msg"]
		payload = fmt.Sprintf("[%s] %s", level, msg)
	case "goroutine":
		nr := int(p.Data["nr"].(float64))
		if nr > 200 {
			payload = fmt.Sprintf("excessive number of active threads: %d", nr)
		}
		c.lock.Lock()
		c.status.ThreadNr = nr
		c.lock.Unlock()
	default:
		return "", ""
	}

	return code.(string), payload
}
