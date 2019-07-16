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
		c.lock.Lock()
		payload = "new block validated: "

		if round, ok := p.Data["round"]; ok {
			c.status.Round = uint64(round.(float64))
			payload = fmt.Sprintf("%s round: %d", payload, c.status.Round)
		}
		if hash, ok := p.Data["blockHash"]; ok {
			c.status.BlockHash = hash.(string)
			payload = fmt.Sprintf("%s hash: %s", payload, c.status.BlockHash)
		}
		if time, ok := p.Data["blockTime"]; ok {
			c.status.BlockTime = time.(float64)
			payload = fmt.Sprintf("%s block time: %.2f", payload, c.status.BlockTime)
		}
		c.lock.Unlock()

	case "warn":
		if level, ok := p.Data["level"]; ok {
			payload = fmt.Sprintf("[%s]", level)
		}
		msg := p.Data["msg"]
		payload = fmt.Sprintf("%s %s", payload, msg)
	case "goroutine":
		if nr, ok := p.Data["nr"]; ok {
			n := int(nr.(float64))
			if n > 200 {
				payload = fmt.Sprintf("excessive number of active threads: %d", nr)
			}
			c.lock.Lock()
			c.status.ThreadNr = n
			c.lock.Unlock()
			break
		}
		return "", ""
	default:
		return "", ""
	}

	return code.(string), payload
}
