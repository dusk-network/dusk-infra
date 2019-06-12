package latency_test

import (
	"bytes"
	"testing"
	"time"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/latency"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestDelay(t *testing.T) {
	w := new(bytes.Buffer)
	m := &monitor.Param{
		Timestamp: time.Now(),
	}

	mntr := latency.New("178.62.193.89")
	assert.NoError(t, mntr.Monitor(w, m))
	assert.NotNil(t, m.Value)
}
