package fake_roundtrip

import (
	"net/http"
)

func NewFakeRoundTripAgent() *FakeRoundTripAgent {
	return &FakeRoundTripAgent{roundTrips: make(map[string]http.RoundTripper)}
}

type FakeRoundTripAgent struct {
	client *FakeRoundTripClient
	roundTrips map[string]http.RoundTripper
}

func (f FakeRoundTripAgent) RoundTrip(r *http.Request) (*http.Response, error) {
//	fmt.Println("agent intercept: ", r.URL)
	if roundTrip := f.roundTrips[f.getKey(*r)]; roundTrip != nil {
		return roundTrip.RoundTrip(r)
	}

	return FourOFour(), nil
}

func (f FakeRoundTripAgent) add(url, method string, roundTrip FakeRoundTrip) {
	roundTrip.fakeRoundTripAgent = &f

	key := f.makeKey(url, method)
	f.roundTrips[key] = roundTrip
	//	fmt.Println("cache len: ", len(f.roundTrips), key)
}

func (f FakeRoundTripAgent) makeKey(url, method string) string {
	return url + ":" + method
}

func (f FakeRoundTripAgent) getKey(r http.Request) string {
	return r.URL.String() + ":" + r.Method
}
