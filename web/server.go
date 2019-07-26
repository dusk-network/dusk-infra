package web

import (
	"fmt"
	"net/http"
	"sync"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/aggregator"

	lg "github.com/sirupsen/logrus"

	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/mux"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Returns true for everything right now
		return true
	},
}

var log = lg.WithField("process", "server")

// Srv is represents the monitoring server. It accepts HTTP and Websocket
// connections and forwards monitor's initial state and data
type Srv struct {
	Monitors []monitor.Mon
	lock     sync.Mutex
	muxConn  *mux.Writer
}

// New creates a Srv by passing the monitors and the aggregator Client
func New(m []monitor.Mon, a *aggregator.Client) *Srv {
	muxConn := mux.New()
	if a != nil {
		muxConn.Add(a)
	}
	return &Srv{
		Monitors: m,
		muxConn:  muxConn,
	}
}

// Serve the monitoring page and upgrade the route `/ws` to websockets listening to the streams of monitoring information
func (s *Srv) Serve(addr string) error {
	for _, mon := range s.Monitors {
		go mon.Wire(s.muxConn)
	}

	d := http.Dir("client")
	fs := http.FileServer(d)
	http.HandleFunc("/stats", s.stats)
	http.Handle("/", fs)

	log.Debugln(fmt.Sprintf("Listening on %s\n", addr))
	return http.ListenAndServe(addr, nil)
}

func (s *Srv) stats(w http.ResponseWriter, r *http.Request) {
	log := log.WithField("api", "stats")
	log.Debugln("beginning upgrade")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Errorln("problem in upgrading the websocket")
		return
	}

	s.lock.Lock()
	id := s.muxConn.Add(c)

	// communicating the initial state to the new connection
	for _, mon := range s.Monitors {
		// this is to force sampling immediately when a new connection arrives
		// so to not show a blank chart
		log.Debugln("forcing sampling on monitors allowing that")
		force, okForce := mon.(monitor.ForceMon)
		if okForce {
			force.ForceSampling()
		}

		log.Debugln("checking stateful monitor")
		init, ok := mon.(monitor.StatefulMon)
		if ok {
			log.WithField("monitor", fmt.Sprintf("%s", mon)).Debugln("sending initial state")
			if err := init.InitialState(s.muxConn); err != nil {
				log.WithError(err).Errorln("problem in initializing the process")
				continue
			}
		}
	}
	s.lock.Unlock()

	defer s.dispose(id, c)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.WithError(err).Debugln("closing websocket")
			} else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithError(err).Warnln("unexpected closing error")
			}
			break
		}
		log.WithField("message", message).Debugln("message received")
	}
}

func (s *Srv) dispose(id uint32, c *websocket.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.muxConn.Remove(id)
	c.Close()
}
