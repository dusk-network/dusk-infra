package mem

import (
	"encoding/json"
	"fmt"
	"io"

	memproc "github.com/shirou/gopsutil/mem"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

type Mem struct{}

func (me *Mem) Monitor(w io.Writer, m *monitor.Param) error {
	v, err := memproc.VirtualMemory()
	if err != nil {
		return err
	}
	m.Value = fmt.Sprintf("%d", int(v.UsedPercent))
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
