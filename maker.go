package fakepoint

import (
	"net/http"
)

func NewFakepointMaker() *FakepointMaker {
	agent := NewAgent()
	maker := &FakepointMaker{
		agent: agent,
	}

	agent.maker = maker
	return maker
}

type FakepointMaker struct {
	agent *Agent
}

func (f FakepointMaker) Client() *http.Client {
	return &http.Client{ Transport: f.agent }
}

func (f FakepointMaker) AddTrip(method, url string, statusCode int) *Fakepoint {
	fr := &Fakepoint{
		method:     method,
		url:        url,
		statusCode: statusCode,
		header:     &http.Header{},
		agent:      f.agent,
	}

	f.agent.add(url, method, fr)
	return fr
}

func (f FakepointMaker) NewGet(url string, statusCode int) *Fakepoint {
	return f.AddTrip("GET", url, statusCode)
}

func (f FakepointMaker) NewPost(url string, statusCode int) *Fakepoint {
	return f.AddTrip("POST", url, statusCode)
}

func (f FakepointMaker) NewPut(url string, statusCode int) *Fakepoint {
	return f.AddTrip("PUT", url, statusCode)
}

func (f FakepointMaker) NewDelete(url string, statusCode int) *Fakepoint {
	return f.AddTrip("DELETE", url, statusCode)
}
