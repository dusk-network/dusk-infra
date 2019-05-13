package collectors

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// websockets in golang do not support concurrent write, so we have to use a Mutex
var mutex = &sync.Mutex{}

type messg struct {
	Timestamp time.Time `json:"timestamp"`
	Metric    string    `json:"metric"`
	Value     string    `json:"value"`
}

// CPUReader keeps reading CPU load values, and transmits them to the websocket channel every n seconds.
func CPUReader(c *websocket.Conn, interval time.Duration) {

	for {
		//spew.Dump(c.conn)
		time.Sleep(interval * time.Second)
		cpuPct, _ := cpu.Percent(0, false)
		err := sendMessage(c, "cpu", fmt.Sprintf("%f", cpuPct))
		if err != nil {
			fmt.Println(err)
			fmt.Println("stopping cpureader...")
			return
		}
	}
}

// DiskReader provides stats about system disk access
func DiskReader(c *websocket.Conn, interval time.Duration) {
	type diskStats struct {
		Mountpoint string  `json:"mountpoint"`
		Percent    float64 `json:"percent"`
		Total      string  `json:"total"`
		Free       string  `json:"free"`
		Used       string  `json:"used"`
	}

	for {
		time.Sleep(interval * time.Second)

		parts, err := disk.Partitions(false)
		if err != nil {
			fmt.Println(err)
		}

		var partitionStats []diskStats

		for _, part := range parts {
			u, err := disk.Usage(part.Mountpoint)
			if err != nil {
				fmt.Println(err)
			}
			d := &diskStats{
				Mountpoint: u.Path,
				Percent:    u.UsedPercent,
				Total:      strconv.FormatUint(u.Total/1024/1024/1024, 10) + " GiB",
				Free:       strconv.FormatUint(u.Free/1024/1024/1024, 10) + " GiB",
				Used:       strconv.FormatUint(u.Used/1024/1024/1024, 10) + " GiB",
			}

			partitionStats = append(partitionStats, *d)

		}
		jsn, _ := json.Marshal(partitionStats)
		sendMessage(c, "dsk", string(jsn))
	}
}

// MemReader provides stats about system memory status
func MemReader(c *websocket.Conn, interval time.Duration) {
	for {
		time.Sleep(interval * time.Second)
		v, _ := mem.VirtualMemory()
		err := sendMessage(c, "mem", v.String())
		if err != nil {
			fmt.Println(err)
			fmt.Println("stopping memreader...")

			return
		}
	}
}

// // NetReader provides stats about network usage
// func NetReader(c *websocket.Conn, interval time.Duration) {
// 	for {
// 		time.Sleep(interval * time.Second)
// 		v, _ := net.NetIOCountersStat()
// 		err := sendMessage(c, "mem", v.String())
// 		if err != nil {
// 			fmt.Println(err)
// 			fmt.Println("stopping memreader...")

// 			return
// 		}
// 	}
// }

func sendMessage(c *websocket.Conn, t string, m string) error {
	mutex.Lock()
	defer mutex.Unlock()

	msg := &messg{
		Metric:    t,
		Value:     m,
		Timestamp: time.Now(),
	}

	j, _ := json.Marshal(msg)

	err := c.WriteJSON(string(j))
	if err != nil {
		return err
	}
	return nil
}
