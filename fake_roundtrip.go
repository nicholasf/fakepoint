package fake_client

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

//1. a factory for producing fake clients
//- should take HTTP Method
//- expected body
//- URL (as a literal or a regex)
//- expected response status code
//- the expected URL parameters with corresponding assertions available via regex (can I get this to work for post vars)
//- perhaps convenience functions for working with JSON posts?

func NewFakeClient(method string, url string, document string) *http.Client {
	fr := &FakeRoundTripper{
		method:     method,
		url:        url,
		document:   document,
		statusCode: 200,
	}

	return &http.Client{Transport: fr}
}

type FakeRoundTripper struct {
	method     string
	url        string
	document   string
	statusCode int
}

type FakeReadCloser struct {
	io.Reader
}

func (f FakeReadCloser) Close() error {
	return nil
}

func (f FakeRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var statusCode int = f.statusCode

	fr := &FakeReadCloser{
		Reader: strings.NewReader(f.document),
	}

	expectedURL, _ := url.Parse(f.url)
	if !urlMatches(*expectedURL, *r.URL) {
		statusCode = 404
	}

	resp := &http.Response{
		Body:       fr,
		StatusCode: statusCode,
	}

	return resp, nil
}

func urlMatches(actual url.URL, expected url.URL) bool {
	match := (actual.Scheme == expected.Scheme) && (actual.Host == expected.Host) && (actual.Path == expected.Path) && (actual.RawQuery == expected.RawQuery)
	if match {
		return true
	}

	return false
}
