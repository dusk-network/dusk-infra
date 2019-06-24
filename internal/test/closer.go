package test

import "io"

type Closer struct {
	io.Writer
	Closed bool
}

func NewCloser(w io.Writer) *Closer {
	return &Closer{
		Writer: w,
		Closed: false,
	}
}

func (c *Closer) Close() error {
	c.Closed = true
	return nil
}
