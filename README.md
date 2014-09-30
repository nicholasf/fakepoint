## Fakepoint

Fakepoint - Fake Endpoint

Create Fake endpoints for HTTP testing. Specify the response data sent back.

[Godocs](http://godoc.org/github.com/nicholasf/fakepoint)

## Install

`go get github.com/nicholasf/fakepoint`

## Example

Use it to simulate a third party API in your tests. 

```golang
maker := NewFakepointMaker()
maker.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }").SetHeader("Content-Type", "application/json")
client := maker.Client()
resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
text, _ := ioutil.ReadAll(resp.Body)
fmt.Println(text) // "{ \"code\": 200 }"
fmt.Println(resp.Header.Get("Content-Type")) //"application/json"
```

## Explanation

First, set up the FakepointMaker.

```golang
maker := NewFakepointMaker()
```

From here, set up an endpoint with `maker.PlanGet`, `maker.PlanPost`, `maker.PlanPut`, or `maker.PlanDelete`. Also specify the http status code you expect, and the final text document.

```golang
maker.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }")
```

You can chain further calls to set headers and increase the frequency of the endpoint:

``` golang
trip := maker.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }")
trip.SetHeader("Content-Type", "application/json").Duplicate(1)
```

Fakepoints only have a lifetime of one request *unless* `Duplicate` is used to specify additional call lifetimes.

This call sets the Content-Type in the response header, and increases the amount of times the client will field this request by 1.

Finally, get the http.Client from the FakepointMaker:

```golang
client := maker.Client()
```

Golang's HTTP design makes setting up fake URLs very easy. Fakepoint simply sets up a `http.Roundtripper` implementation that resolves HTTP requests against a map of fake endpoints. Voila!

Fakepoint will cleanly handle 302s by following redirects, and supplying a 'new-location' resource to automate the follows.

Mostly written in Santa Monica, L.A., 2014, while working for Rockpool Labs.

##FAQs

"What's the benefit of fakepoint over http://golang.org/pkg/net/http/httptest?"

The [httptest](http://golang.org/pkg/net/http/httptest/) package provides an easy to set up HTTP server for testing. 

I wanted to specify particular URLs, and to write tests to assert that they were called. 

I've used tools like [nock](https://github.com/pgte/nock), [fakeweb](https://github.com/chrisk/fakeweb), etc. in other languages and I wanted something similar.

In particular, I see myself writing a lot of integration code against 3rd party HTTP services, and I want a nice DSL in my tests to prove my code works. All I have to do is expose the call to the client and it's suddenly very testable:

```golang
package notifications_test

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"github.com/nicholasf/fakepoint"
	"allocation/notifications"
)

func TestOpsgenie(t *testing.T) {
	Convey("An Error log is sent to the notifier", t, func() {
		maker := fakepoint.NewFakepointMaker()
		maker.PlanPost("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }")
		resp, err := notifications.Requester(*maker.Client(), ("https://api.opsgenie.com/v1/json/alert"), []byte(``))
		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, 200)
	})
}


```

"Can I do transfer-coding = "chunked" and get nice responses?"

Not yet, but that'd be grand. I might get around to it if I ever need it.

## License

Super open.