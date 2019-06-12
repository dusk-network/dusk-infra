package cpu_test

import (
	"bytes"
	"testing"
	"time"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/cpu"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestCPU(t *testing.T) {
	w := new(bytes.Buffer)
	m := &monitor.Param{
		Timestamp: time.Now(),
	}

	mntr := &cpu.CPU{}
	assert.NoError(t, mntr.Monitor(w, m))
	assert.NotNil(t, m.Value)
}
