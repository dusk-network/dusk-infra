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
		if time != nil {
			switch time.(type) {
			case string:
				c.status.BlockTime = time.(string)
			case float64:
				c.status.BlockTime = fmt.Sprintf("%.2f", time.(float64))
			}
		}
		c.status.Round = uint64(round.(float64))
		c.lock.Unlock()

		payload = fmt.Sprintf("new block validated: round %d, hash %s, block time %sms", c.status.Round, hash, time)
	case "warn":
		level := p.Data["level"]
		msg := p.Data["msg"]
		payload = fmt.Sprintf("[%s] %s", level, msg)
	}

	return code.(string), payload
}
