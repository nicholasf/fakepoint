package fake_roundtrip

import (
	"net/http"
	"net/url"
	"fmt"
)

const REDIRECTED_LOCATION = "/new-location/"

type FakeRoundTrip struct {
	statusCode int
	method     string
	url        string
	document   string
	header 	*http.Header
	fakeRoundTripAgent *FakeRoundTripAgent
}

func (f FakeRoundTrip) SetStatusCode(code int) *FakeRoundTrip {
	f.statusCode = code

	locationRequired := code == 302 || code == 201 || code == 202

	if locationRequired {
		setDefaultLocationHeader(&f)
	}

	return &f
}

func (f FakeRoundTrip) SetHeader(key string, value string) *FakeRoundTrip {
	f.header.Set(key, value)
	return &f
}

func (f FakeRoundTrip) SetURL(url string) *FakeRoundTrip {
	f.url = url
	return &f
}

func (f FakeRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	var statusCode int = f.statusCode

//	fmt.Println("URL1: ", r.URL, statusCode)
	if !f.requestMatches(*r) {
//		fmt.Println("URL2: ", r.URL, statusCode)

		statusCode = 404
		return FourOFour(), nil
	}

//	if (statusCode == 302) {
//		redirectedURL := r.URL.Scheme + "://" + r.URL.Host + REDIRECTED_LOCATION
//		fmt.Println("URL reset to ", f.url, " for next req")
//		fmt.Println("\n\n\n")

	resp := &http.Response {
		Body:       NewFakeReadCloser(f.document),
		StatusCode: statusCode,
		Header: *f.header,
	}

	return resp, nil
}

func (f FakeRoundTrip) requestMatches(r http.Request) bool {
	if (f.method != r.Method) {
		fmt.Println("Method ", f.method, " does not match ", r.Method)
		return false
	}

	actual, _ := url.Parse(f.url)
	expected := r.URL

	pathMatch := actual.Path == expected.Path
//	fmt.Println("pathMatch: ", pathMatch)
	schemeMatch := actual.Scheme == expected.Scheme
//	fmt.Println("schemeMatch: ", schemeMatch)
	hostMatch := actual.Host == expected.Host
//	fmt.Println("hostMatch: ", hostMatch)
	rawQueryMatch := actual.RawQuery == expected.RawQuery
//	fmt.Println("rawQueryMatch: ", rawQueryMatch)
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
	fmt.Println("****** 404")
	resp := &http.Response {
		Body:    NewFakeReadCloser("Unknown"),
		StatusCode: 404,
	}

	return resp
}
