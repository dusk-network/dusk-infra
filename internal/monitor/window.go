package monitor

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// DataWindow is a shifting window of Data
type DataWindow []map[string]interface{}

// NewDataWindow create a new DataWindow
func NewDataWindow() DataWindow {
	return make([]map[string]interface{}, 0)
}

// Add all elements of a DataWindow
func (d DataWindow) Add(other DataWindow) DataWindow {
	return append(d, other...)
}

// Append data to DataWindow
func (d DataWindow) Append(m map[string]interface{}) DataWindow {
	now := time.Now()
	if _, ok := m["timestamp"]; !ok {
		m["timestamp"] = now.Format(time.RFC3339)
	}

	d = append(d, m)
	sort.Sort(&d)

	if len(d) > 0 {
		// removing obsolete data
		for _, prev := range d {
			t := prev["timestamp"].(string)
			tr, _ := time.Parse(time.RFC3339, t)

			if now.Sub(tr) > MaxTimeSpan {
				d = d[1:]
				continue
			}
			break
		}
	}

	return d
}

// Swap is part of sort.Interface
func (d *DataWindow) Swap(i, j int) {
	(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
}

// Less is part of sort.Interface
func (d *DataWindow) Less(i, j int) bool {
	ti, tj := (*d)[i]["timestamp"].(string), (*d)[j]["timestamp"].(string)
	timeI, _ := time.Parse(time.RFC3339, ti)
	timeJ, _ := time.Parse(time.RFC3339, tj)
	return timeI.Before(timeJ)
}

// Len is part of sort.Interface
func (d *DataWindow) Len() int {
	return len(*d)
}

// Datum is a timestamped value
type Datum struct {
	T   time.Time
	Val float64
}

// MarshalJSON takes care of correctly formatting data in RFC3339 standard
func (d Datum) MarshalJSON() ([]byte, error) {
	stamp := d.T.Format(time.RFC3339)
	str := fmt.Sprintf("{ \"value\":%.2f,\"timestamp\":\"%s\"}", d.Val, stamp)
	return []byte(str), nil
}

// UnmarshalJSON takes care of correctly parsing data from RFC3339 standard
func (d *Datum) UnmarshalJSON(b []byte) error {
	var err error
	md := make(map[string]interface{})
	if err = json.Unmarshal(b, &md); err != nil {
		return err
	}
	d.Val = md["value"].(float64)
	t := md["timestamp"]
	d.T, err = time.Parse(time.RFC3339, t.(string))
	return err
}

// Window is a utility struct to help with calculation
type Window []Datum

// MaxTimeSpan sets the time for when data gets obsolete
var MaxTimeSpan = 5 * time.Minute

// NewWindow is the constructor for a shifting window
func NewWindow() Window {
	return make(Window, 0)
}

// CalculateAvg of this Window
func (w Window) CalculateAvg() float64 {
	if len(w) == 0 {
		return 0
	}

	var sum float64
	for _, d := range w {
		sum += d.Val
	}

	return sum / float64(len(w))
}

// Add is a convenience method to append all Datum at once to the Window
func (w Window) Add(d Window) Window {
	return append(w, d...)
}

// Append to the Window. This method does not mutate the Window
func (w Window) Append(d ...float64) Window {
	now := time.Now()
	if len(w) > 0 {
		// removing obsolete data
		for _, prev := range w {
			if now.Sub(prev.T) > MaxTimeSpan {
				w = w[1:]
				continue
			}
			break
		}
	}

	return append(w, newData(d...)...)
}

func newData(ds ...float64) []Datum {
	now := time.Now()
	data := make([]Datum, len(ds))
	for i := range ds {
		data[i] = Datum{now, ds[i]}
	}
	return data
}
