package disk

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/shirou/gopsutil/disk"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

// Disk is a disk monitor
type Disk struct {
	monitor.Window
}

type stats struct {
	Mountpoint string  `json:"mountpoint"`
	Percent    float64 `json:"percent"`
	Total      string  `json:"total"`
	Free       string  `json:"free"`
	Used       string  `json:"used"`
}

// New creates a new Disk monitor
func New() *Disk {
	return &Disk{
		Window: make(monitor.Window, 0),
	}
}

// Monitor writes the disk space monitor on the writer and saves the result on a shifting window
func (d *Disk) Monitor(w io.Writer, m *monitor.Param) error {
	parts, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	for _, part := range parts {
		if part.Mountpoint == "/" {
			u, err := disk.Usage(part.Mountpoint)
			if err != nil {
				return err
			}

			s := &stats{
				Mountpoint: u.Path,
				Percent:    u.UsedPercent,
				Total:      fmt.Sprintf("%s GiB", strconv.FormatUint(u.Total/1024/1024/1024, 10)),
				Free:       fmt.Sprintf("%s GiB", strconv.FormatUint(u.Free/1024/1024/1024, 10)),
				Used:       fmt.Sprintf("%s GiB", strconv.FormatUint(u.Used/1024/1024/1024, 10)),
			}

			m.Value = fmt.Sprintf("%.2f", s.Percent)
			d.Window = d.Append(s.Percent)

			b, err := json.Marshal(m)
			if err != nil {
				return err
			}
			if _, err := w.Write(b); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
