## Fake-Endpoint

Hands you a http.Client that can mock HTTP endpoints, letting you specify the response data returned. 

## Install

`go get github.com/nicholasf/go-fake-endpoint`

## Example

Use it to simulate a third party API in your tests. 

```golang
client := NewFakeClient()
client.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }").SetHeader("Content-Type", "application/json")
resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
text, _ := ioutil.ReadAll(resp.Body)
fmt.Println(text) // "{ \"code\": 200 }"
fmt.Println(resp.Header.Get("Content-Type")) //"application/json"
```

## Explanation

First, set up the client.

```golang
client := NewFakeClient()
```

From here, set up an endpoint with `client.PlanGet`, `client.PlanPost`, `client.PlanPut`, or `client.PlanDlete` with the matching call on the client. Also specify the http status code you expect, and the final text document.

```golang
client.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }")
```

You can chain further calls to set headers and increase the frequency of the endpoint:

``` golang
trip := client.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }")
trip.SetHeader("Content-Type", "application/json").Duplicate(1)
```

This call sets the Content-Type in the response header, and increases the amount of times the client will field this request by 1.
