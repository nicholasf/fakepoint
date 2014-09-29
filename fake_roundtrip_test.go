package fake_roundtrip

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

func TestFakeRoundTrip(t *testing.T) {
	Convey("Basic DSL - planning roundtrips", t, func() {
		client := NewFakeClient()

		Convey("it returns the document with a 200", func() {
			client.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "hello world").SetHeader("Content-Type", "text/plain")

			resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
			text, _ := ioutil.ReadAll(resp.Body)
			So(string(text), ShouldEqual, "hello world")
			So(resp.StatusCode, ShouldEqual, 200)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "text/plain")
		})

		Convey("it doesnt resolve to the wrong URL", func() {
			resp, _ := client.Get("https://somethingelse.com")
			So(resp.StatusCode, ShouldEqual, 404)
		})

		Convey("the header", func() {
			client.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "hello world").SetHeader("Content-Type", "application/json")
			resp, err := client.Get("https://api.opsgenie.com/v1/json/alert")
			So(err, ShouldBeNil)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "application/json")
		})
	})

	Convey("The Client handles multiple fake round trips", t, func() {
		client := NewFakeClient()
		Convey("distinct URLs", func() {
			client.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "hello world")
			client.PlanGet("https://another.system.com", 200, "not with a bang but a whimper")
			resp1, err := client.Get("https://api.opsgenie.com/v1/json/alert")
			text, _ := ioutil.ReadAll(resp1.Body)
			So(err, ShouldBeNil)
			So(string(text), ShouldEqual, "hello world")
			So(resp1.StatusCode, ShouldEqual, 200)

			resp2, err2 := client.Get("https://another.system.com")
			text, _ = ioutil.ReadAll(resp2.Body)
			So(err2, ShouldBeNil)
			So(string(text), ShouldEqual, "not with a bang but a whimper")
			So(resp2.StatusCode, ShouldEqual, 200)
		})

		Convey("Identical URLs, different HTTP methods", func() {
			client.PlanPost("abc.com/greeting", 201, "hello world1")
			client.PlanGet("abc.com/greeting", 200, "hello world2")

			resp1, err := client.Post("abc.com/greeting", "application/text", bytes.NewReader([]byte{'g', 'o'}))
			text, _ := ioutil.ReadAll(resp1.Body)
			So(err, ShouldBeNil)
			So(string(text), ShouldEqual, "hello world1")
			So(resp1.StatusCode, ShouldEqual, 201)

			resp2, err2 := client.Get("abc.com/greeting")
			text, _ = ioutil.ReadAll(resp2.Body)
			So(err2, ShouldBeNil)
			So(string(text), ShouldEqual, "hello world2")
			So(resp2.StatusCode, ShouldEqual, 200)
		})
	})

	Convey("302s (following redirects)", t, func() {
		client := NewFakeClient()

		Convey("302 generates an interal 200 for /new-location/ according to 302 rules", func() {
			client.PlanGet("abc.com", 302, "")
			resp, err := client.Get("abc.com")
			So(err, ShouldBeNil)
			So(resp.Header.Get("Location"), ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("the status code, with Location headers for 201, 202", func() {
			codes := []int{201, 202}

			for _, code := range codes {
				client.PlanGet("abc.com", code, "")
				resp, err := client.Get("abc.com")
				So(err, ShouldBeNil)
				So(resp.Header.Get("Location"), ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, code)
			}
		})
	})

	Convey("Repetition of Endpoint default to a repetition of 1", t, func() {
		client := NewFakeClient()

		Convey("defaults to 1", func() {
			client.PlanGet("abc.com", 200, "")
			resp, err := client.Get("abc.com")
			resp2, err := client.Get("abc.com")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			So(resp2.StatusCode, ShouldEqual, 404)
		})

		Convey("is configurable", func() {
			client.PlanGet("abc.com", 200, "").Duplicate(1)
			resp, err := client.Get("abc.com")
			resp2, err := client.Get("abc.com")
			resp3, err := client.Get("abc.com")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			So(resp2.StatusCode, ShouldEqual, 200)
			So(resp3.StatusCode, ShouldEqual, 404)
		})
	})
}
