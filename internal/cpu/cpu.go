package cpu

import (
	"encoding/json"
	"fmt"
	"io"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"

	cpumon "github.com/shirou/gopsutil/cpu"
)

// CPU monitoring process
type CPU struct {
	data monitor.Window
}

// New create a new CPU Monitoring process
func New() *CPU {
	return &CPU{
		data: make(monitor.Window, 0),
	}
}

// Monitor writes the current value of the CPU on the writer
func (cp *CPU) Monitor(w io.Writer, m *monitor.Param) error {
	cpuPct, _ := cpumon.Percent(0, false)
	m.Value = fmt.Sprintf("%f", cpuPct[0])
	cp.data = cp.data.Append(cpuPct[0])
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
