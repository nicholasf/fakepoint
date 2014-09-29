package fake_roundtrip

import (
	"net/http"
	"net/url"
)

const REDIRECTED_LOCATION = "/new-location/"

type FakeRoundTrip struct {
	statusCode  int
	method      string
	url         string
	document    string
	header      *http.Header
	agent       *FakeRoundTripAgent
	repetitions int
}

func (f *FakeRoundTrip) SetStatusCode(code int) *FakeRoundTrip {
	f.statusCode = code

	locationRequired := code == 302 || code == 201 || code == 202

	if locationRequired {
		setDefaultLocationHeader(f)
	}

	return f
}

func (f *FakeRoundTrip) SetHeader(key string, value string) *FakeRoundTrip {
	f.header.Set(key, value)
	return f
}

func (f *FakeRoundTrip) SetURL(url string) *FakeRoundTrip {
	f.url = url
	return f
}

func (f *FakeRoundTrip) Duplicate(num int) *FakeRoundTrip {
	f.agent.increaseCount(f.url, f.method, num)
	return f
}

func (f FakeRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	var statusCode int = f.statusCode

	if !f.requestMatches(*r) {
		statusCode = 404
		return FourOFour(), nil
	}

	if statusCode == 302 {
		redirectURL := REDIRECTED_LOCATION
		f.agent.client.AddTrip(f.method, redirectURL, 200, "")
		setDefaultLocationHeader(&f)
	}

	resp := &http.Response{
		Body:       NewFakeReadCloser(f.document),
		StatusCode: statusCode,
		Header:     *f.header,
	}

	return resp, nil
}

func (f FakeRoundTrip) requestMatches(r http.Request) bool {
	if f.method != r.Method {
		return false
	}

	actual, _ := url.Parse(f.url)
	expected := r.URL

	pathMatch := actual.Path == expected.Path
	schemeMatch := actual.Scheme == expected.Scheme
	hostMatch := actual.Host == expected.Host
	rawQueryMatch := actual.RawQuery == expected.RawQuery
	match := pathMatch && schemeMatch && hostMatch && rawQueryMatch

	if match {
		return true
	}

	return false
}

//per handle location headers: http://en.wikipedia.org/wiki/HTTP_location
func setDefaultLocationHeader(f *FakeRoundTrip) {
	f.SetHeader("Location", REDIRECTED_LOCATION)
}

func FourOFour() *http.Response {
	resp := &http.Response{
		Body:       NewFakeReadCloser("Unknown"),
		StatusCode: 404,
	}

	return resp
}
