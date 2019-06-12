package collectors

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
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
		err := sendMessage(c, "cpu", fmt.Sprintf("%f", cpuPct[0]))
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

		for _, part := range parts {
			if part.Mountpoint == "/" {
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

				err = sendMessage(c, "dsk", fmt.Sprintf("%d", int(d.Percent)))
				if err != nil {
					fmt.Println(err)
					fmt.Println("stopping diskreader...")
					return
				}
			}
		}
	}
}

// MemReader provides stats about system memory status
func MemReader(c *websocket.Conn, interval time.Duration) {
	for {
		time.Sleep(interval * time.Second)
		v, _ := mem.VirtualMemory()
		err := sendMessage(c, "mem", fmt.Sprintf("%d", int(v.UsedPercent)))
		if err != nil {
			fmt.Println(err)
			fmt.Println("stopping memreader...")
			return
		}
	}
}

// LatencyReader pings the Dusk voucher seeder to estimate network latency
func LatencyReader(c *websocket.Conn, interval time.Duration) {
	for {
		time.Sleep(interval * time.Second)
		p := func(addr string) (float64, error) {
			_, dur, err := Ping(addr)
			if err != nil {
				fmt.Println(err)
				return 0, err
			}
			return (float64(dur) / 1000000), nil
		}

		// Pings the voucher seeder
		latency, _ := p("178.62.193.89")

		err := sendMessage(c, "net", fmt.Sprintf("%f", latency))
		if err != nil {
			fmt.Println(err)
			fmt.Println("stopping latencyreader...")
			return
		}
	}
}

// LogReader monitors the Dusk log file and sends updates to the dashboard
func LogReader(c *websocket.Conn, logfile string) {
	file, err := os.Open(logfile)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	var fileBuffLines = make([]string, 0)

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		fileBuffLines = append(fileBuffLines, string(line))
	}

	length := len(fileBuffLines)

	// Send last 10 lines
	lineCount := 10

	if lineCount > length {
		lineCount = length
	}

	for i := length - lineCount; i < length; i++ {
		err = sendMessage(c, "log", fileBuffLines[i])
	}

	fi, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}

	t, err := tail.TailFile(logfile, tail.Config{Follow: true, Poll: true, Location: &tail.SeekInfo{fi.Size(), os.SEEK_SET}})
	if err != nil {
		fmt.Println(err)
	}

	for line := range t.Lines {
		err = sendMessage(c, "log", line.Text)
		if err != nil {
			fmt.Println(err)
			fmt.Println("stopping logreader...")
			return
		}
	}
}

func sendMessage(c *websocket.Conn, t string, m string) error {
	fmt.Println("=> Sending message..." + t)
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
