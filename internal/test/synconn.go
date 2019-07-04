package test

import (
	"bytes"
	"encoding/json"
)

type JsonReadWriter struct {
	*bytes.Buffer
	Closed bool
}

func NewWriter() *JsonReadWriter {
	return &JsonReadWriter{
		Buffer: new(bytes.Buffer),
	}
}

func (s *JsonReadWriter) WriteJSON(v interface{}) error {

	switch v.(type) {
	case string:
		if _, err := s.Write([]byte(v.(string))); err != nil {
			return err
		}
	default:
		d, err := json.Marshal(v)
		if err != nil {
			return err
		}

		if _, err := s.Write([]byte(d)); err != nil {
			return err
		}
	}

	return nil
}

func (s *JsonReadWriter) ReadMessage() (int, []byte, error) {
	b := s.Bytes()
	return len(b), b, nil
}

func (s *JsonReadWriter) Close() error {
	s.Closed = true
	return nil
}
