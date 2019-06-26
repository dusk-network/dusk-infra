package monitor_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestMonitor(t *testing.T) {
	r, w := io.Pipe()
	buf := new(bytes.Buffer)
	s := mockSupervisor{t: "test"}

	m := monitor.New(s, 10*time.Millisecond, "test")
	go m.Wire(w)

	//giving enough time to monitor.Wire
	time.Sleep(time.Millisecond)

	// writing 3 times the test string
	b := make([]byte, 4)
	for i := 0; i < 3; i++ {
		_, err := r.Read(b)
		if !assert.NoError(t, err) {
			return
		}
		buf.Write(b)
	}

	assert.Equal(t, "testtesttest", buf.String())
}

type mockSupervisor struct {
	t string
}

func (m mockSupervisor) Monitor(w io.Writer, p *monitor.Param) error {
	p.Value = m.t
	if _, err := w.Write([]byte(p.Value)); err != nil {
		return err
	}
	return nil
}
