package logstream_test

import (
	"bytes"
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

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const file = "/tmp/dusk.soc"

var unixSoc = fmt.Sprintf("unix://%s", file)

func init() {
	log.SetLevel(log.TraceLevel)
}

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

	packet := `{ "code": "round", "blockTime": 30, "round": "30", "blockHash": "12345" }`
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
		assert.Equal(t, float64(30), p.Data["blockTime"])
		assert.Equal(t, "30", p.Data["round"])
		assert.Equal(t, "12345", p.Data["blockHash"])
	}

	if !assert.NoError(t, nl.InitialState(synconn)) {
		t.FailNow()
	}

	buf := new(bytes.Buffer)

	if !assert.NoError(t, nl.InitialState(buf)) {
		t.FailNow()
	}

	p := &monitor.Param{}
	if !assert.NoError(t, json.Unmarshal(buf.Bytes(), p)) {
		t.FailNow()
	}
	assert.Equal(t, "30", p.Data["round"])
	win := p.Data["blockTimes"].([]interface{})
	v := win[0].(map[string]interface{})
	assert.Equal(t, float64(30), v["value"].(float64))
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

func await(t *testing.T, nl *logstream.Monitor, d time.Duration) {
	//giving enough time to the server to start
	select {
	case err := <-nl.ErrChan:
		fmt.Println(err)
		assert.FailNow(t, "unexpected error")
	case <-time.After(d):
		return
	}
}
