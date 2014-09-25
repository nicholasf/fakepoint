package fake_round_trip

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"fmt"
)

const REDIRECTED_LOCATION = "/new-location/"

func NewFakeClient(method string, url string, document string) *FakeRoundTripClient {
	fr := &FakeRoundTrip {
		method:     method,
		url:        url,
		document:   document,
		statusCode: 200,
		header: &http.Header{},
	}

	fakeClient := &FakeRoundTripClient{ fakeRoundTrip: fr }
	fakeClient.Transport = fr
	return fakeClient
}

type FakeRoundTripClient struct {
	http.Client
	fakeRoundTrip *FakeRoundTrip
}

func (f FakeRoundTripClient) SetStatusCode(code int) *FakeRoundTripClient {
	f.fakeRoundTrip.statusCode = code

	locationRequired := code == 302 || code == 201 || code == 202

	if locationRequired {
		setDefaultLocationHeader(&f)
	}

	return &f
}

func (f FakeRoundTripClient) SetResponseHeader(key string, value string) *FakeRoundTripClient {
	f.fakeRoundTrip.header.Set(key, value)
	return &f
}

func (f FakeRoundTripClient) SetURL(url string) *FakeRoundTripClient {
	f.fakeRoundTrip.url = url
	return &f
}

type FakeRoundTrip struct {
	statusCode int
	method     string
	url        string
	document   string
	header 	*http.Header
}

type FakeReadCloser struct {
	io.Reader
}

func (f FakeReadCloser) Close() error {
	return nil
}

func (f FakeRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Println("*************** ", f.url)
	var statusCode int = f.statusCode

	fr := &FakeReadCloser{
		Reader: strings.NewReader(f.document),
	}

	expectedURL, _ := url.Parse(f.url)
//	fmt.Println("expectedURL: ", expectedURL)
//	fmt.Println("actualURL: ",  *r.URL)
	if !urlMatches(*r.URL, *expectedURL) {
		statusCode = 404
	}

	if (statusCode == 302) {
		f.url = r.URL.Scheme + "://" + r.URL.Host + REDIRECTED_LOCATION
		fmt.Println("URL reset to ", f.url, " for next req")
		fmt.Println("\n\n\n")
	}

	expectedURL, _ = url.Parse(f.url)
	fmt.Println("expectedURL2: ", expectedURL)

	resp := &http.Response {
		Body:       fr,
		StatusCode: statusCode,
		Header: *f.header,
	}

	return resp, nil
}

func urlMatches(actual url.URL, expected url.URL) bool {
	match := (actual.Scheme == expected.Scheme) && (actual.Host == expected.Host) && (actual.Path == expected.Path) && (actual.RawQuery == expected.RawQuery)
	fmt.Println("1: ", (actual.Scheme == expected.Scheme))
	fmt.Println("2: ", (actual.Host == expected.Host))
	fmt.Println("3: ", (actual.Path == expected.Path))
	fmt.Println("actual.Path ", actual.Path)
	fmt.Println("expected.Path ", expected.Path)

	fmt.Println("4: ", (actual.RawQuery == expected.RawQuery))

	if match {
		return true
	}

	return false
}

//per handle location headers: http://en.wikipedia.org/wiki/HTTP_location
func setDefaultLocationHeader(f *FakeRoundTripClient) {
	f.SetResponseHeader("Location", REDIRECTED_LOCATION)
}
