package test

import (
	"bytes"
	"errors"
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
	// by, err := json.Marshal(v)
	bys, ok := v.(string)
	if !ok {
		return errors.New("Please use string")
	}

	if _, err := s.Write([]byte(bys)); err != nil {
		return err
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
