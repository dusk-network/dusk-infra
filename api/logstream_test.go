package logstream_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	logstream "gitlab.dusk.network/dusk-core/node-monitor/api"
	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/test"

	"github.com/stretchr/testify/assert"
)

const file = "/tmp/dusk.soc"

var unixSoc = fmt.Sprintf("unix://%s", file)

func TestPipe(t *testing.T) {
	_ = os.Remove(file)
	defer os.Remove(file)

	h, err := url.Parse(unixSoc)
	if err != nil {
		panic(err)
	}
	nl := logstream.New(h)
	w := test.NewWriter()
	synconn := j.New(w)

	go nl.Wire(synconn)

	await(t, nl, 1*time.Millisecond)

	packet := `{ "Pippo": "pluto" }`
	if assert.NoError(t, send(packet)) {
		await(t, nl, 1*time.Millisecond)
		p := monitor.NewParam("log")
		_, b, err := synconn.ReadMessage()
		if !assert.NoError(t, err) {
			t.FailNow()
		}

		if !assert.NoError(t, json.Unmarshal(b, p)) {
			t.FailNow()
		}
		assert.Equal(t, "pluto", p.Data["Pippo"])
	}
}

func send(data string) error {
	u, _ := url.Parse(unixSoc)
	c, err := net.Dial(u.Scheme, u.Path)
	if err != nil {
		return err
	}
	if _, err := c.Write([]byte(data)); err != nil {
		return err
	}
	return nil

}

func await(t *testing.T, nl *logstream.LogStreamMonitor, d time.Duration) {
	//giving enough time to the server to start
	select {
	case err := <-nl.ErrChan:
		assert.FailNow(t, "%s\n", err)
	case <-time.After(d):
		return
	}
}
