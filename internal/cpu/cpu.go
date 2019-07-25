package cpu

import (
	"io"

	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"

	cpumon "github.com/shirou/gopsutil/cpu"
)

var _ monitor.StatefulSampler = (*CPU)(nil)

// CPU monitoring process
type CPU struct {
	monitor.Window
}

// String returns the name of this monitor
func (c *CPU) String() string {
	return "cpu"
}

// New create a new CPU Monitoring process
func New() *CPU {
	return &CPU{
		Window: make(monitor.Window, 0),
	}
}

// Monitor writes the current value of the CPU on the writer
func (c *CPU) Monitor(w io.Writer, m *monitor.Param) error {
	cpuPct, err := cpumon.Percent(0, false)
	if err != nil {
		return err
	}
	m.Window = m.Window.Append(cpuPct[0])
	c.Window = c.Add(m.Window)
	if err := j.Write(w, m); err != nil {
		return err
	}
	return nil
}

// InitialState as defined in the StatefulMon interface
func (c *CPU) InitialState(w io.Writer) error {
	if len(c.Window) > 0 {
		m := monitor.NewParam("cpu")
		m.Window = c.Window
		if err := j.Write(w, m); err != nil {
			return err
		}
	}
	return nil
}
