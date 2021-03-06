package aggregator_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/aggregator"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func newP(m string, v float64, d map[string]interface{}) *monitor.Param {

	w := monitor.NewWindow()
	w = w.Append(v)
	return &monitor.Param{
		Metric:    m,
		Window:    w,
		Data:      d,
		Timestamp: time.Now(),
	}
}

var tt = []struct {
	p        *monitor.Param
	expected string
}{
	{
		newP("cpu", 94.02, nil),
		"high CPU load (94.02%)",
	},
	{
		newP("disk", 95.11, nil),
		"little or no Disk space left (95.11%)",
	},
	{
		newP("mem", 95.11, nil),
		"high memory usage (95.11%)",
	},
	{
		newP("latency", 20, map[string]interface{}{
			"error": "my balls itch",
		}),
		"my balls itch",
	},
	{
		newP("latency", 200, nil),
		"network too slow. Latency more than 150ms (200ms)",
	},
	{
		newP("log", 0, map[string]interface{}{
			"code":      "round",
			"blockHash": "pippo",
			"round":     uint64(4),
			"blockTime": float64(4),
		}),
		"new block validated: round 4, hash pippo, block time 4.00ms",
	},
	{
		newP("log", 0, map[string]interface{}{
			"code":  "warn",
			"error": "pippo",
			"msg":   "pluto",
			"level": "titanic",
		}),
		"[titanic] pluto",
	}}

func TestClient(t *testing.T) {
	for _, tData := range tt {
		fn := hndl(t, tData.expected)
		fireSrv(t, tData.p, fn)
	}
}

func parseUrl(u string, t *testing.T) *url.URL {
	uri, err := url.Parse(u)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return uri
}

func fireSrv(t *testing.T, param *monitor.Param, fn func(http.ResponseWriter, *http.Request)) {
	srv := httptest.NewTLSServer(http.HandlerFunc(fn))
	defer srv.Close()

	uri := parseUrl(srv.URL, t)
	sri := "http://localhost:8080"
	a := aggregator.New(uri, sri, "12345", "x.x.x.x", "pippo")

	b, err := json.Marshal(param)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	_, err = a.Write(b)
	assert.NoError(t, err)
}

func hndl(t *testing.T, expected string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		alert := &aggregator.Alert{}
		if !assert.NotEmpty(t, req.Body) {
			t.FailNow()
		}

		if !assert.NoError(t, json.NewDecoder(req.Body).Decode(&alert)) {
			t.FailNow()
		}

		if !assert.Equal(t, expected, alert.Content) {
			t.FailNow()
		}

		if !assert.NotEmpty(t, alert.Ipv4) {
			t.FailNow()
		}
		if !assert.NotEmpty(t, alert.Hostname) {
			t.FailNow()
		}
		_, err := rw.Write([]byte("OK"))

		if !assert.NoError(t, err) {
			t.FailNow()
		}
	}
}
