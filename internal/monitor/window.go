package monitor

import (
	"encoding/json"
	"fmt"
	"time"
)

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
var MaxTimeSpan = 1 * time.Minute

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
