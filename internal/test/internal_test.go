package test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/cpu"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/disk"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/mem"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var monitors = []struct {
	name string
	mon  monitor.Supervisor
}{
	{name: "cpu", mon: &cpu.CPU{}},
	{name: "disk", mon: &disk.Disk{}},
	{name: "mem", mon: &mem.Mem{}},
	// {name: "latency", mon: latency.New("178.62.193.89")},
}

func TestSuite(t *testing.T) {
	for _, tt := range monitors {
		w := new(bytes.Buffer)
		m := &monitor.Param{
			Timestamp: time.Now(),
		}

		if assert.NoError(t, tt.mon.Monitor(w, m), "error messages for process %s", tt.name) {
			assert.NotNilf(t, m.Value, tt.name)
		}
	}
}
