package monitor

import "time"

type datum struct {
	tm  time.Time
	val float64
}

// Window is a utility struct to help with calculation
type Window []datum

// MaxTimeSpan sets the time for when data gets obsolete
var MaxTimeSpan = 1 * time.Minute

// CalculateAvg of this Window
func (w Window) CalculateAvg() float64 {
	if len(w) == 0 {
		return 0
	}

	var sum float64
	for _, d := range w {
		sum += d.val
	}

	return sum / float64(len(w))
}

// Append to the Window. This method does not mutate the Window
func (w Window) Append(d float64) Window {
	now := time.Now()
	if len(w) > 0 {
		// removing obsolete data
		for _, prev := range w {
			if now.Sub(prev.tm) > MaxTimeSpan {
				w = w[1:]
				continue
			}
			break
		}

	}
	return append(w, datum{val: d, tm: now})
}
