package cpu_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/cpu"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestCPU(t *testing.T) {
	w := new(bytes.Buffer)
	m := monitor.NewParam("cpu")
	mntr := cpu.New()
	assert.NoError(t, mntr.Monitor(w, m))
	assert.Equal(t, 1, len(m.Window))
}

func TestInitialState(t *testing.T) {
	c := cpu.New()
	c.Window = c.Append(2.1, 3.2, 4.3)
	w := new(bytes.Buffer)
	if !assert.NoError(t, c.InitialState(w)) {
		t.FailNow()
	}

	p := monitor.NewParam("test")
	if !assert.NoError(t, json.Unmarshal(w.Bytes(), p)) {
		t.FailNow()
	}

	assert.Equal(t, 3, len(p.Window))
	assert.Equal(t, 2.1, p.Window[0].Val)
}
