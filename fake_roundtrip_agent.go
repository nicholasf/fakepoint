package fake_roundtrip

import (
	"net/http"
)

func NewFakeRoundTripAgent() *FakeRoundTripAgent {
	return &FakeRoundTripAgent{
		roundTrips:      make(map[string]http.RoundTripper),
		roundTripsCount: make(map[string]int),
	}
}

type FakeRoundTripAgent struct {
	client          *FakeRoundTripClient
	roundTrips      map[string]http.RoundTripper
	roundTripsCount map[string]int
}

func (f FakeRoundTripAgent) RoundTrip(r *http.Request) (*http.Response, error) {
	key := f.getKey(*r)
	if roundTrip := f.roundTrips[key]; roundTrip != nil {
		resp, err := roundTrip.RoundTrip(r)
		f.resolveCount(roundTrip.(FakeRoundTrip))
		return resp, err
	}

	return FourOFour(), nil
}

func (f FakeRoundTripAgent) add(url, method string, roundTrip FakeRoundTrip) {
	roundTrip.agent = &f

	key := f.makeKey(url, method)
	f.roundTrips[key] = roundTrip
	f.roundTripsCount[key] = 1
}

func (f FakeRoundTripAgent) increaseCount(url, method string, num int) {
	key := f.makeKey(url, method)
	f.roundTripsCount[key] += num
}

func (f FakeRoundTripAgent) resolveCount(roundTrip FakeRoundTrip) {
	key := f.makeKey(roundTrip.url, roundTrip.method)
	f.roundTripsCount[key]--
	if f.roundTripsCount[key] == 0 {
		f.remove(roundTrip)
	}
}

func (f FakeRoundTripAgent) remove(roundTrip FakeRoundTrip) {
	key := f.makeKey(roundTrip.url, roundTrip.method)
	delete(f.roundTrips, key)
	delete(f.roundTripsCount, key)
}

func (f FakeRoundTripAgent) makeKey(url, method string) string {
	return url + ":" + method
}

func (f FakeRoundTripAgent) getKey(r http.Request) string {
	return r.URL.String() + ":" + r.Method
}
