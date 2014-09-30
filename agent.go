package fakepoint

import (
	"net/http"
)

func NewAgent() *Agent {
	return &Agent{
		fakepoints:      make(map[string]http.RoundTripper),
		fakepointsCount: make(map[string]int),
	}
}

type Agent struct {
	maker          *FakepointMaker
	fakepoints      map[string]http.RoundTripper
	fakepointsCount map[string]int
}

func (f Agent) RoundTrip(r *http.Request) (*http.Response, error) {
	key := f.getKey(*r)
	if roundTrip := f.fakepoints[key]; roundTrip != nil {
		resp, err := roundTrip.RoundTrip(r)
		f.resolveCount(roundTrip.(Fakepoint))
		return resp, err
	}

	return FourOFour(), nil
}

func (f Agent) add(url, method string, roundTrip Fakepoint) {
	roundTrip.agent = &f

	key := f.makeKey(url, method)
	f.fakepoints[key] = roundTrip
	f.fakepointsCount[key] = 1
}

func (f Agent) increaseCount(url, method string, num int) {
	key := f.makeKey(url, method)
	f.fakepointsCount[key] += num
}

func (f Agent) resolveCount(roundTrip Fakepoint) {
	key := f.makeKey(roundTrip.url, roundTrip.method)
	f.fakepointsCount[key]--
	if f.fakepointsCount[key] == 0 {
		f.remove(roundTrip)
	}
}

func (f Agent) remove(roundTrip Fakepoint) {
	key := f.makeKey(roundTrip.url, roundTrip.method)
	delete(f.fakepoints, key)
	delete(f.fakepointsCount, key)
}

func (f Agent) makeKey(url, method string) string {
	return url + ":" + method
}

func (f Agent) getKey(r http.Request) string {
	return r.URL.String() + ":" + r.Method
}
