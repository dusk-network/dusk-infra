package aggregator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	lg "github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/ip"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

var log = lg.WithField("process", "aggregator")

type Client struct {
	uri        *url.URL
	httpclient *http.Client
	token      string
}

func New(uri *url.URL, token string) *Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &Client{
		uri:        uri,
		httpclient: client,
		token:      token,
	}
}

func (c *Client) WriteJSON(v interface{}) error {
	var payload string
	p := v.(*monitor.Param)

	switch p.Metric {
	case "latency":
		payload = c.serializeLatency(p)
	case "disk":
		payload = c.serializeDisk(p)
	case "cpu":
		payload = c.serializeCpu(p)
	case "log":
		payload = c.serializeLog(p)
	case "mem":
		payload = c.serializeMem(p)
	case "tail":
		payload = ""
	default:
		log.WithField("metric", p.Metric).Warnln("unrecognized metric. Can't forward to aggregator")
		return nil
	}

	if len(payload) > 0 {
		c.send(payload)
	}

	return nil
}

// Write locally on param that gets serialized into an HTTP call toward the aggregator
func (c *Client) Write(b []byte) (int, error) {
	p := &monitor.Param{}
	if err := json.Unmarshal(b, p); err != nil {
		// we are unmarshalling something that has been marshalled by this monitoring program
		//if we have problems it means the program is wrong
		panic(err)
	}
	_ = c.WriteJSON(p)
	return len(b), nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) ReadMessage() (int, []byte, error) {
	return 0, []byte{}, nil
}

func (c *Client) send(payload string) {
	var ipv4, hostname string
	var err error
	if ipv4, err = ip.Retrieve(); err != nil {
		log.WithError(err).Warnln("cannot retrieve the IP of the machine")
	}
	hostname, err = os.Hostname()
	if err != nil {
		log.WithError(err).Warnln("cannot retrieve the hostname of the machine")
	}
	alert := &Alert{
		Content:  payload,
		Hostname: hostname,
		Ipv4:     ipv4,
	}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(alert)

	req, _ := http.NewRequest("POST", c.uri.String(), b)
	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-type", "application/json; charset=utf-8")
	res, err := c.httpclient.Do(req)
	if err != nil {
		log.WithError(err).Warnln("problems in posting to the aggregator")
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Warnln("problems in reading the response")
	}
	res.Body.Close()
}

type Alert struct {
	Content  string `json:"content"`
	Ipv4     string `json:"ipv4"`
	Hostname string `json:"hostname"`
}
