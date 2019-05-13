package collectors

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type messg struct {
	Timestamp time.Time `json:"timestamp"`
	Metric    string    `json:"metric"`
	Value     string    `json:"value"`
}

type response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

// CPUReader keeps reading CPU load values, and transmits them to the websocket channel every n seconds.
func CPUReader(c *websocket.Conn, interval time.Duration) {

	for {
		time.Sleep(interval * time.Second)
		cpuPct, _ := cpu.Percent(0, false)
		msg := &messg{
			Metric:    "cpu",
			Value:     fmt.Sprintf("%f", cpuPct),
			Timestamp: time.Now(),
		}
		j, _ := json.Marshal(msg)
		spew.Dump(msg)
		fmt.Println(string(j))

		err := c.WriteJSON(string(j))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// DiskReader provides stats about system disk access
func DiskReader(c *net.Conn) {
}

// MemReader provides stats about system memory status
func MemReader(c *net.Conn) {
	for {
		v, _ := mem.VirtualMemory()
	}
}
