package web

import (
	"fmt"
	"net/http"
	"sync"

	lg "github.com/sirupsen/logrus"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
	"gitlab.dusk.network/dusk-core/node-monitor/web/json"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Returns true for everything right now
		return true
	},
}

var log = lg.WithField("process", "server")

type (
	Srv struct {
		Monitors []monitor.Mon
	}

	synConn struct {
		*websocket.Conn
		lock sync.RWMutex
	}
)

func (s *synConn) WriteJSON(v string) error {
	log.WithField("payload", v).Debugln("sending outgoing packet")
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Conn.WriteJSON(v)
}

func (s *synConn) ReadMessage() (int, []byte, error) {
	log.Debugln("reading incoming messages")
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Conn.ReadMessage()
}

// Serve the monitoring page and upgrade the route `/ws` to websockets listening to the streams of monitoring information
func (s *Srv) Serve(addr string) error {

	d := http.Dir("static")
	fs := http.FileServer(d)
	http.HandleFunc("/stats", s.stats)
	http.Handle("/", fs)

	log.Debugln(fmt.Sprintf("Listening on %s\n", addr))
	return http.ListenAndServe(addr, nil)
}

func (s *Srv) stats(w http.ResponseWriter, r *http.Request) {
	log := log.WithField("api", "stats")
	log.Infoln("beginning upgrade")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Errorln("problem in upgrading the websocket")
		return
	}
	defer s.dispose(c)

	sc := &synConn{c, sync.RWMutex{}}
	for _, mon := range s.Monitors {
		go json.New(sc, mon).Connect()
	}

	quitC := make(chan struct{})
	<-quitC

	// for {
	// 	_, _, err := sc.ReadMessage()
	// 	if err != nil {
	// 		log.WithError(err).Errorln("problem in receiving messages")
	// 		return
	// 	}
	// 	log.Debugln("got mail!")
	// }
}

func (s *Srv) dispose(c *websocket.Conn) {
	// this supposedly triggers an error on the reader part of the infinite for loop of the Bridge.Connect method
	defer c.Close()
	for _, mon := range s.Monitors {
		mon.Disconnect()
	}
}
