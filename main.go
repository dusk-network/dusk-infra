package main

import (
	"flag"
	"log"
	"net/http"

	collectors "./collectors"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Returns true for everything right now
		return true
	},
}

func main() {

	flag.Parse()
	log.SetFlags(0)
	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/stats", stats)
	http.Handle("/", fs)

	//go consumer(messages)
	log.Fatal(http.ListenAndServe(*addr, nil))

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
	go collectors.DiskReader(c, 5)

	// Hold on Listen Channel
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}
