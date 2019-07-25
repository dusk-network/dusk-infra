package mem

import (
	"io"

	memproc "github.com/shirou/gopsutil/mem"
	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

// Mem is the general memory monitoring process
type Mem struct {
	monitor.Window
}

// New creates a new Mem process
func New() *Mem {
	return &Mem{
		Window: make(monitor.Window, 0),
	}
}

// Monitor writes the mem reading to a Writer. It also saves a value on the shifting Window
func (me *Mem) Monitor(w io.Writer, m *monitor.Param) error {
	v, err := memproc.VirtualMemory()
	if err != nil {
		return err
	}
	used := v.UsedPercent
	m.Window = m.Window.Append(used)
	me.Window = me.Window.Add(m.Window)
	if err := j.Write(w, m); err != nil {
		return err
	}
	return nil
}

// InitialState as defined by the StatefulMon interface
func (me *Mem) InitialState(w io.Writer) error {
	p := monitor.NewParam("mem")
	p.Window = me.Window
	if err := j.Write(w, p); err != nil {
		return err
	}
	return nil
}

// String returns the name of this monitor
func (me *Mem) String() string {
	return "mem"
}
