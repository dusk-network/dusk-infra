package disk_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/disk"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestDisk(t *testing.T) {
	w := new(bytes.Buffer)
	m := &monitor.Param{
		Timestamp: time.Now(),
	}

	mntr := disk.New()
	assert.NoError(t, mntr.Monitor(w, m))
	assert.Equal(t, 1, len(m.Window))
	assert.Greater(t, m.Window[0].Val, 1.0)
}

func TestInitialState(t *testing.T) {
	d := disk.New()
	d.Window = d.Append(2.1, 3.2, 4.3)
	w := new(bytes.Buffer)
	if !assert.NoError(t, d.InitialState(w)) {
		t.FailNow()
	}

	p := monitor.NewParam("disk")
	if !assert.NoError(t, json.Unmarshal(w.Bytes(), p)) {
		t.FailNow()
	}

	assert.Equal(t, 3, len(p.Window))
	assert.Equal(t, 2.1, p.Window[0].Val)
}
