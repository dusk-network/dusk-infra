package web

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"

	lg "github.com/sirupsen/logrus"

	j "gitlab.dusk.network/dusk-core/node-monitor/internal/json"
	"gitlab.dusk.network/dusk-core/node-monitor/internal/monitor"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Returns true for everything right now
		return true
	},
}

var log = lg.WithField("process", "server")

type WriterMux struct {
	sync.Mutex
	io.Writer
	conns map[uint32]*j.SynConn
}

func NewWriterMux() *WriterMux {
	return &WriterMux{
		conns: make(map[uint32]*j.SynConn),
	}
}

func (mux *WriterMux) Write(b []byte) (int, error) {
	if mux.Writer != nil {
		return mux.Writer.Write(b)
	}
	log.Debugln("No writer specified yet. Dropping packet")
	return 0, nil
}

func (mux *WriterMux) Add(jw j.JsonReadWriter) uint32 {
	mux.Lock()
	defer mux.Unlock()
	id := rand.Uint32()
	mux.conns[id] = j.New(jw)
	mux.rebuildWriter()
	return id
}

func (mux *WriterMux) Remove(id uint32) {
	mux.Lock()
	defer mux.Unlock()
	delete(mux.conns, id)
	mux.rebuildWriter()
}

func (mux *WriterMux) rebuildWriter() {
	var curConn []io.Writer
	for _, conn := range mux.conns {
		curConn = append(curConn, conn)
	}

	mux.Writer = io.MultiWriter(curConn...)
}

type Srv struct {
	Monitors []monitor.Mon
	lock     sync.Mutex
	muxConn  *WriterMux
}

// Serve the monitoring page and upgrade the route `/ws` to websockets listening to the streams of monitoring information
func (s *Srv) Serve(addr string) error {
	s.muxConn = NewWriterMux()
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
	log.Infoln("beginning upgrade")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Errorln("problem in upgrading the websocket")
		return
	}

	s.lock.Lock()
	id := s.muxConn.Add(c)
	s.lock.Unlock()

	defer s.dispose(id, c)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithError(err).Debugln("closing websocket")
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
