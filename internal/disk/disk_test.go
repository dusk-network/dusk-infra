package disk_test

import (
	"bytes"
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

	mntr := &disk.Disk{}
	assert.NoError(t, mntr.Monitor(w, m))
	assert.NotNil(t, m.Value)
}
