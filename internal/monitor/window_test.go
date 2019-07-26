package monitor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestDataWindow(t *testing.T) {
	w := monitor.NewDataWindow()

	m1 := fillMap(time.Now().Add(2 * time.Minute))
	m2 := fillMap(time.Now())
	w = w.Append(m1)
	w = w.Append(m2)

	assert.Equal(t, m2, w[0])
	assert.Equal(t, m1, w[1])

	pastTime := time.Now().Add(-1 * (monitor.MaxTimeSpan + 1*time.Minute))
	past := fillMap(pastTime)
	w = w.Append(past)
	assert.Equal(t, 2, len(w))
}

func fillMap(clock time.Time) map[string]interface{} {
	m := make(map[string]interface{})
	m["timestamp"] = clock.Format(time.RFC3339)
	return m
}
