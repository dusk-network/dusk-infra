package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strconv"

	collectors "./collectors"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var logfile = "/Users/fulvio/log.txt"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Returns true for everything right now
		return true
	},
}

func main() {

	isRoot := checkPrivileges()
	if isRoot == false {
		log.Fatal("Error: must be run as root!")
	}

	flag.Parse()
	log.SetFlags(0)
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/stats", stats)
	http.Handle("/", fs)

	//go consumer(messages)
	log.Fatal(http.ListenAndServe(*addr, nil))

}

func checkPrivileges() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}
	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		log.Fatal(err)
	}

	if i == 0 {
		return true
	}
	return false
}

func stats(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// Spawn CPU parser
	go collectors.CPUReader(c, 5)
	go collectors.MemReader(c, 8)
	go collectors.LatencyReader(c, 10)
	go collectors.DiskReader(c, 5)
	// go collectors.LogReader(c, logfile)

	// Hold on Listen Channel
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}
