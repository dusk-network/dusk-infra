package aggregator

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
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
	lock       sync.RWMutex
	status     *Status
}

func New(uri *url.URL, token string) *Client {
	var err error
	var hostname, ipv4 string
	// tr := &http.Transport{
	//     TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	client := &http.Client{
		// Transport: tr,
		Timeout: 10 * time.Second,
	}

	ipv4, err = ip.Retrieve()
	if err != nil {
		panic(err)
	}

	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	hc := &Client{
		uri:        uri,
		httpclient: client,
		token:      token,
		status: &Status{
			Ipv4:     ipv4,
			Hostname: hostname,
		},
	}
	go hc.sendUpdate()
	return hc
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
	alert := &Alert{
		Content:  payload,
		Hostname: c.status.Hostname,
		Ipv4:     c.status.Ipv4,
	}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(alert)
	c.forward(b, "alert")
}

func (c *Client) forward(b *bytes.Buffer, endpoint string) {
	req, _ := http.NewRequest("POST", c.uri.String()+endpoint, b)
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

func (c *Client) sendUpdate() {
	tick := time.NewTicker(20 * time.Second)
	for {
		<-tick.C
		b := new(bytes.Buffer)
		c.lock.RLock()
		_ = json.NewEncoder(b).Encode(c.status)
		c.lock.RUnlock()
		c.forward(b, "update")
	}
}

type Alert struct {
	Content  string `json:"content"`
	Ipv4     string `json:"ipv4"`
	Hostname string `json:"hostname"`
}

type Status struct {
	Ipv4      string  `json:"ipv4"`
	Hostname  string  `json:"hostname"`
	CPU       float32 `json:"cpu"`
	Disk      float32 `json:"disk"`
	Round     uint64  `json:"height"`
	BlockTime string  `json:"blockTime"`
	BlockHash string  `json:"blockHash"`
	Latency   float32 `json:"latency"`
	Mem       float32 `json:"mem"`
}
