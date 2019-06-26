package io

import (
	"sync"

	lg "github.com/sirupsen/logrus"
)

var log = lg.WithField("process", "SynConn")

type MessageReader interface {
	ReadMessage() (int, []byte, error)
}

type JsonWriter interface {
	WriteJSON(v interface{}) error
	Close() error
}

type JsonReadWriter interface {
	MessageReader
	JsonWriter
}

type SynConn struct {
	sync.RWMutex
	JsonReadWriter
}

func New(w JsonReadWriter) *SynConn {
	return &SynConn{
		JsonReadWriter: w,
	}
}

func (s *SynConn) Write(b []byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	if err := s.WriteJSON(string(b)); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (s *SynConn) ReadMessage() (int, []byte, error) {
	log.Debugln("reading incoming messages")
	s.RLock()
	defer s.RUnlock()
	return s.JsonReadWriter.ReadMessage()
}
