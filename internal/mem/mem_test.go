package mem_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/mem"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestMem(t *testing.T) {
	w := new(bytes.Buffer)
	m := monitor.NewParam("mem")
	mntr := mem.New()

	assert.NoError(t, mntr.Monitor(w, m))
	assert.Equal(t, 1, len(m.Window))
}

func TestInitialState(t *testing.T) {
	m := mem.New()
	m.Window = m.Append(2.1, 3.2, 4.3)

	w := new(bytes.Buffer)
	if !assert.NoError(t, m.InitialState(w)) {
		t.FailNow()
	}

	p := monitor.NewParam("test")
	if !assert.NoError(t, json.Unmarshal(w.Bytes(), p)) {
		t.FailNow()
	}
	assert.Equal(t, 3, len(p.Window))
	assert.Equal(t, 2.1, p.Window[0].Val)
}
