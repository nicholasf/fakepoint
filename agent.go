package fakepoint

import (
	"net/http"
)

func NewAgent() *Agent {
	return &Agent{
		fakepoints:      make(map[string]*Fakepoint),
		fakepointsCount: make(map[string]int),
	}
}

type Agent struct {
	maker          *FakepointMaker
	fakepoints      map[string]*Fakepoint
	fakepointsCount map[string]int
}

func (a Agent) RoundTrip(r *http.Request) (*http.Response, error) {
	key := a.getKey(*r)
	if fakepoint := a.fakepoints[key]; fakepoint != nil {
		roundTrip := http.RoundTripper(*fakepoint)
		resp, err := roundTrip.RoundTrip(r)
		a.resolveCount(roundTrip.(Fakepoint))
		return resp, err
	}

	return FourOFour(), nil
}

func (a Agent) add(url, method string, roundTrip *Fakepoint) {
	roundTrip.agent = &a

	key := a.makeKey(url, method)
	a.fakepoints[key] = roundTrip
	a.fakepointsCount[key] = 1
}

func (a Agent) increaseCount(url, method string, num int) {
	key := a.makeKey(url, method)
	a.fakepointsCount[key] += num
}

func (a Agent) resolveCount(roundTrip Fakepoint) {
	key := a.makeKey(roundTrip.url, roundTrip.method)
	a.fakepointsCount[key]--
	if a.fakepointsCount[key] == 0 {
		a.remove(roundTrip)
	}
}

func (a Agent) remove(roundTrip Fakepoint) {
	key := a.makeKey(roundTrip.url, roundTrip.method)
	delete(a.fakepoints, key)
	delete(a.fakepointsCount, key)
}

func (a Agent) makeKey(url, method string) string {
	return url + ":" + method
}

func (a Agent) getKey(r http.Request) string {
	return r.URL.String() + ":" + r.Method
}
