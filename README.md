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
client := maker.client()
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
client := maker.client()
```

Golang's HTTP design makes setting up fake URLs very easy. Fakepoint simply sets up a `http.Roundtripper` implementation that resolves HTTP requests against a map of fake endpoints. Voila!

Fakepoint will cleanly handle 302s by following redirects, and supplying a 'new-location' resource to automate the follows.

Mostly written in Santa Monica, 2014, while working for Rockpool Labs.

## License

Super open.