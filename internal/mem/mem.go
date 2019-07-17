package mem

import (
	"encoding/json"
	"fmt"
	"io"

	memproc "github.com/shirou/gopsutil/mem"
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
	m.Value = fmt.Sprintf("%.2f", v.UsedPercent)
	me.Window = me.Append(v.UsedPercent)

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
