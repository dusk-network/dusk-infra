package main

import (
	"bytes"
	"fmt"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/latency"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
)

func main() {
	l := latency.New("178.62.193.89")
	buf := new(bytes.Buffer)

	p := &monitor.Param{}
	if err := l.Monitor(buf, p); err != nil {
		fmt.Println("oops")
		panic(err)
	}
	fmt.Println(buf.String())
}
