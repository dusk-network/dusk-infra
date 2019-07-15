package aggregator

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
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
var tolerance = time.Minute * 3

// Client to the aggregator of alerts and status. It connects to a URL and send JSON packets
type Client struct {
	uri        *url.URL
	httpclient *http.Client
	token      string
	lock       sync.RWMutex
	status     *Status
	alerts     map[string]*Alert
}

// New creates a new Aggregator Client and sets up the connection
func New(uri *url.URL, srv string, token string) *Client {
	var err error
	var hostname, ipv4 string
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	if uri.Scheme == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
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
			Srv:      srv,
		},
		alerts: make(map[string]*Alert),
	}
	go hc.sendUpdate()
	return hc
}

// WriteJSON accepts a JSON encodable struct, checks if an alert needs to be sent to the aggregator and updates the status
func (c *Client) WriteJSON(v interface{}) error {
	var payload, code string
	p := v.(*monitor.Param)

	switch p.Metric {
	case "latency":
		payload = c.serializeLatency(p)
	case "disk":
		payload = c.serializeDisk(p)
	case "cpu":
		payload = c.serializeCpu(p)
	case "log":
		code, payload = c.serializeLog(p)
		if code == "" {
			return nil
		}
	case "mem":
		payload = c.serializeMem(p)
	case "tail":
		payload = ""
	default:
		log.WithField("metric", p.Metric).Warnln("unrecognized metric. Can't forward to aggregator")
		return nil
	}

	if len(payload) > 0 {
		if formerAlert, ok := c.alerts[code]; ok {
			span := time.Since(formerAlert.createdAt)
			if span < tolerance {
				return nil
			}
		}
		// we do not send new block messages as alerts
		// the information is already sent through a Status update
		if code != "round" {
			a := c.send(payload)
			c.alerts[code] = a
		}
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

// Close as defined by the json.JsonReadWriter interface
func (c *Client) Close() error {
	return nil
}

// ReadMessage as defined by the json.JsonReadWriter interface
func (c *Client) ReadMessage() (int, []byte, error) {
	return 0, []byte{}, nil
}

func (c *Client) send(payload string) *Alert {
	alert := &Alert{
		Content:   payload,
		Hostname:  c.status.Hostname,
		Ipv4:      c.status.Ipv4,
		createdAt: time.Now(),
	}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(alert)
	c.forward(b, "alert")
	return alert
}

func (c *Client) forward(r io.Reader, endpoint string) {
	tgt := c.uri.String() + "/" + endpoint
	logger := new(bytes.Buffer)
	tr := io.TeeReader(r, logger)
	req, _ := http.NewRequest("POST", tgt, tr)
	req.Header.Add("Authorization", c.token)
	req.Header.Add("Content-type", "application/json; charset=utf-8")
	res, err := c.httpclient.Do(req)
	if err != nil {
		log.WithError(err).Warnln("problems in posting to the aggregator")
		return
	}

	defer res.Body.Close()

	log.WithFields(lg.Fields{
		"payload": logger.String(),
		"url":     tgt,
	}).Debugln("sending " + endpoint)

	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.WithError(err).Warnln("problems in reading the response")
		return
	}

	log.WithField("response", string(resp)).Debugln(endpoint + " sent")
}

func (c *Client) sendUpdate() {
	tick := time.NewTicker(20 * time.Second)
	for {
		<-tick.C
		b := new(bytes.Buffer)
		c.lock.RLock()
		log.WithField("status", c.status).Debugln("sending update")
		_ = json.NewEncoder(b).Encode(c.status)
		c.lock.RUnlock()
		c.forward(b, "update")
	}
}

// Alert is the json encodable struct sent to the aggregator to signal an unexpected alert
// or an error situation
type Alert struct {
	Content   string `json:"content"`
	Ipv4      string `json:"ipv4"`
	Hostname  string `json:"hostname"`
	createdAt time.Time
}

// Status is the json encodable struct sent to the aggregator
type Status struct {
	Ipv4      string  `json:"ipv4"`
	Hostname  string  `json:"hostname"`
	Srv       string  `json:"srv"`
	CPU       float64 `json:"cpu"`
	Disk      float64 `json:"disk"`
	Round     uint64  `json:"height"`
	BlockTime float64 `json:"blockTime"`
	BlockHash string  `json:"blockHash"`
	Latency   float64 `json:"latency"`
	Mem       float64 `json:"memory"`
	ThreadNr  int     `json:"thread"`
}
