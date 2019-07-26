package disk

import (
	"io"

	"github.com/shirou/gopsutil/disk"
	log "github.com/sirupsen/logrus"
	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
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

func (d *Disk) String() string {
	return "disk"
}

// Monitor writes the disk space monitor on the writer and saves the result on a shifting window
func (d *Disk) Monitor(w io.Writer, m *monitor.Param) error {
	parts, err := disk.Partitions(false)
	if err != nil {
		return err
	}

	for _, part := range parts {
		log.WithFields(log.Fields{
			"process":    "disk",
			"mountpoint": part.Mountpoint,
		}).Debugln("reading disk usage")

		if part.Mountpoint == "/" || part.Mountpoint == "/root" {
			u, err := disk.Usage(part.Mountpoint)
			if err != nil {
				return err
			}

			/*
				s := &stats{
					Mountpoint: u.Path,
					Percent:    u.UsedPercent,
					Total:      fmt.Sprintf("%s GiB", strconv.FormatUint(u.Total/1024/1024/1024, 10)),
					Free:       fmt.Sprintf("%s GiB", strconv.FormatUint(u.Free/1024/1024/1024, 10)),
					Used:       fmt.Sprintf("%s GiB", strconv.FormatUint(u.Used/1024/1024/1024, 10)),
				}
			*/

			used := u.UsedPercent
			m.Window = m.Window.Append(used)
			d.Window = d.Window.Add(m.Window)
			if err := j.Write(w, m); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

// InitialState returns the current known window of readings
func (d *Disk) InitialState(w io.Writer) error {
	p := monitor.NewParam("disk")
	p.Window = d.Window
	if err := j.Write(w, p); err != nil {
		return err
	}
	return nil
}
