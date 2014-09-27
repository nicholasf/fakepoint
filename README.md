WIP - nothing to read here yet.

## Fake-Roundtrip

Behaviours for faking endpoint URLs, response headers and bodies, using http.Roundtrippers via the http.Client

## Install

`go get github.com/nicholasf/go-fake-roundtrip`

## Example
```
    client := NewFakeClient()

    Convey("it returns the document with a 200", func() {
        client.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "hello world").SetHeader("Content-Type", "text/plain")

        resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
        text, _ := ioutil.ReadAll(resp.Body)
        So(string(text), ShouldEqual, "hello world")
        So(resp.StatusCode, ShouldEqual, 200)
        So(resp.Header.Get("Content-Type"), ShouldEqual, "text/plain")
    })

```

## Explanation




WIP.

Returns a client containing a fake round trip for mock http testing.
