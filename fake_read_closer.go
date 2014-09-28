package fake_roundtrip

import (
	"io"
	"strings"
)

func NewFakeReadCloser(body string) *FakeReadCloser {
	fr := &FakeReadCloser{
		Reader: strings.NewReader(body),
	}

	return fr
}

type FakeReadCloser struct {
	io.Reader
}

func (f FakeReadCloser) Close() error {
	return nil
}
