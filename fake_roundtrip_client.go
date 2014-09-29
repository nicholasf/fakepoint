package fake_roundtrip

import (
	"net/http"
)

func NewFakeClient() *FakeRoundTripClient {
	agent := NewFakeRoundTripAgent()
	fakeClient := &FakeRoundTripClient{
		agent: agent,
	}

	agent.client = fakeClient
	fakeClient.Transport = agent
	return fakeClient
}

type FakeRoundTripClient struct {
	http.Client
	agent *FakeRoundTripAgent
}

func (f FakeRoundTripClient) AddTrip(method, url string, statusCode int, document string) *FakeRoundTrip {
	fr := &FakeRoundTrip{
		method:     method,
		url:        url,
		document:   document,
		statusCode: statusCode,
		header:     &http.Header{},
		agent:      f.agent,
	}

	f.agent.add(url, method, *fr)
	return fr
}

func (f FakeRoundTripClient) PlanGet(url string, statusCode int, document string) *FakeRoundTrip {
	return f.AddTrip("GET", url, statusCode, document)
}

func (f FakeRoundTripClient) PlanPost(url string, statusCode int, document string) *FakeRoundTrip {
	return f.AddTrip("POST", url, statusCode, document)
}

func (f FakeRoundTripClient) PlanPut(url string, statusCode int, document string) *FakeRoundTrip {
	return f.AddTrip("PUT", url, statusCode, document)
}

func (f FakeRoundTripClient) PlanDelete(url string, statusCode int, document string) *FakeRoundTrip {
	return f.AddTrip("DELETE", url, statusCode, document)
}
