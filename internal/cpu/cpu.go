package cpu

import (
	"encoding/json"
	"fmt"
	"io"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"

	cpumon "github.com/shirou/gopsutil/cpu"
)

type CPU struct{}

func (cp *CPU) Monitor(w io.Writer, m *monitor.Param) error {
	cpuPct, _ := cpumon.Percent(0, false)
	m.Value = fmt.Sprintf("%f", cpuPct[0])
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
