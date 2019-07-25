package test

import (
	"bytes"
	"encoding/json"
)

// JSONReadWriter is a mock of a Websocket connection
type JSONReadWriter struct {
	*bytes.Buffer
	Closed bool
}

// NewWriter creates a new Mock Websocket connection
func NewWriter() *JSONReadWriter {
	return &JSONReadWriter{
		Buffer: new(bytes.Buffer),
	}
}

// WriteJSON mimics the Websocket function
func (s *JSONReadWriter) WriteJSON(v interface{}) error {
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

// ReadMessage mimics the Websocket function
func (s *JSONReadWriter) ReadMessage() (int, []byte, error) {
	b := s.Bytes()
	return len(b), b, nil
}

// Close mimics the Websocket function
func (s *JSONReadWriter) Close() error {
	s.Closed = true
	return nil
}
