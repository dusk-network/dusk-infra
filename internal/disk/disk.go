package disk

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/shirou/gopsutil/disk"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

type Disk struct{}

type stats struct {
	Mountpoint string  `json:"mountpoint"`
	Percent    float64 `json:"percent"`
	Total      string  `json:"total"`
	Free       string  `json:"free"`
	Used       string  `json:"used"`
}

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

			d := &stats{
				Mountpoint: u.Path,
				Percent:    u.UsedPercent,
				Total:      fmt.Sprintf("%s GiB", strconv.FormatUint(u.Total/1024/1024/1024, 10)),
				Free:       fmt.Sprintf("%s GiB", strconv.FormatUint(u.Free/1024/1024/1024, 10)),
				Used:       fmt.Sprintf("%s GiB", strconv.FormatUint(u.Used/1024/1024/1024, 10)),
			}

			m.Value = fmt.Sprintf("%d", int(d.Percent))
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
