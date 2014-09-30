package fakepoint

import (
	"io"
	"strings"
)

func NewReadCloser(body string) *ReadCloser {
	fr := &ReadCloser{
		Reader: strings.NewReader(body),
	}

	return fr
}

type ReadCloser struct {
	io.Reader
}

func (f ReadCloser) Close() error {
	return nil
}
