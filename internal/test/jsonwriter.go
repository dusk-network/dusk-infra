package test

import (
	"encoding/json"
	"sync"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

type Writer struct {
	sync.RWMutex
	p *monitor.Param
}

func NewJsonWriter() *Writer {
	p := &monitor.Param{}
	return &Writer{p: p}
}

func (w *Writer) WriteJSON(s string) error {
	w.Lock()
	defer w.Unlock()
	return json.Unmarshal([]byte(s), &w.p)
}

func (w *Writer) Get() monitor.Param {
	w.RLock()
	defer w.RUnlock()
	return *w.p
}
