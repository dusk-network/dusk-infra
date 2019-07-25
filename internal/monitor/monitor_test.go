package monitor_test

import (
	"encoding/binary"
	"io"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func TestMonitor(t *testing.T) {
	test := float64(3.2)
	r, w := io.Pipe()
	s := mockSampler{t: []float64{test}}

	m := monitor.New(s, 10*time.Millisecond, "test")
	go m.Wire(w)

	//giving enough time to monitor.Wire
	time.Sleep(time.Millisecond)

	// writing 3 times the test value
	var b [8]byte
	for i := 0; i < 3; i++ {
		_, err := r.Read(b[:])
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		// bits
		bits := binary.LittleEndian.Uint64(b[:])
		// bits to float64
		f := math.Float64frombits(bits)

		assert.Equal(t, test, f)
	}
}

type mockSampler struct {
	t []float64
}

func (m mockSampler) Monitor(w io.Writer, p *monitor.Param) error {
	var b [8]byte
	for _, s := range m.t {
		binary.LittleEndian.PutUint64(b[:], math.Float64bits(s))
		if _, err := w.Write(b[:]); err != nil {
			return err
		}
	}
	p.Window = p.Window.Append(m.t...)
	return nil
}
