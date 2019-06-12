package mem_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/mem"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestMem(t *testing.T) {
	w := new(bytes.Buffer)
	m := &monitor.Param{
		Timestamp: time.Now(),
	}

	mntr := &mem.Mem{}
	assert.NoError(t, mntr.Monitor(w, m))
	assert.NotNil(t, m.Value)
}
