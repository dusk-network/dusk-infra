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
			b := c.status.blockTime.Append(time.(float64))
			avg := b.CalculateAvg()
			c.status.blockTime = b
			c.status.BlockTime = avg
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
			n := nr.(float64)

			c.lock.Lock()
			tn := c.status.threads.Append(n)
			avg := tn.CalculateAvg()
			c.status.threads = tn
			c.status.ThreadNr = int(avg)
			c.lock.Unlock()
			if avg > 200 {
				payload = fmt.Sprintf("excessive number of active threads: %d", int(avg))
			}
			break
		}
		return "", ""
	default:
		return "", ""
	}

	return code.(string), payload
}
