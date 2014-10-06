package fakepoint

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
	"net/http"
)

func TestFakeRoundTrip(t *testing.T) {
	Convey("Basic DSL - planning roundtrips", t, func() {
		maker := NewFakepointMaker()

		Convey("it returns the document with a 200", func() {
			maker.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "hello world").SetHeader("Content-Type", "text/plain")

			resp, _ := maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
			text, _ := ioutil.ReadAll(resp.Body)
			So(string(text), ShouldEqual, "hello world")
			So(resp.StatusCode, ShouldEqual, 200)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "text/plain")
		})

		Convey("it doesnt resolve to the wrong URL", func() {
			resp, _ := maker.Client().Get("https://somethingelse.com")
			So(resp.StatusCode, ShouldEqual, 404)
		})

		Convey("the header", func() {
			maker.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "{ \"code\": 200 }").SetHeader("Content-Type", "application/json")
			resp, err := maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
			So(err, ShouldBeNil)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "application/json")
		})
	})

	Convey("The maker handles multiple fake round trips", t, func() {
		maker := NewFakepointMaker()
		Convey("distinct URLs", func() {
			maker.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "hello world")
			maker.PlanGet("https://another.system.com", 200, "not with a bang but a whimper")
			resp1, err := maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
			text, _ := ioutil.ReadAll(resp1.Body)
			So(err, ShouldBeNil)
			So(string(text), ShouldEqual, "hello world")
			So(resp1.StatusCode, ShouldEqual, 200)

			resp2, err2 := maker.Client().Get("https://another.system.com")
			text, _ = ioutil.ReadAll(resp2.Body)
			So(err2, ShouldBeNil)
			So(string(text), ShouldEqual, "not with a bang but a whimper")
			So(resp2.StatusCode, ShouldEqual, 200)
		})

		Convey("Identical URLs, different HTTP methods", func() {
			maker.PlanPost("abc.com/greeting", 201, "hello world1")
			maker.PlanGet("abc.com/greeting", 200, "hello world2")

			resp1, err := maker.Client().Post("abc.com/greeting", "application/text", bytes.NewReader([]byte{'g', 'o'}))
			text, _ := ioutil.ReadAll(resp1.Body)
			So(err, ShouldBeNil)
			So(string(text), ShouldEqual, "hello world1")
			So(resp1.StatusCode, ShouldEqual, 201)

			resp2, err2 := maker.Client().Get("abc.com/greeting")
			text, _ = ioutil.ReadAll(resp2.Body)
			So(err2, ShouldBeNil)
			So(string(text), ShouldEqual, "hello world2")
			So(resp2.StatusCode, ShouldEqual, 200)
		})
	})

	Convey("302s (following redirects)", t, func() {
		maker := NewFakepointMaker()

		Convey("302 generates an interal 200 for /new-location/ according to 302 rules", func() {
			maker.PlanGet("abc.com", 302, "")
			resp, err := maker.Client().Get("abc.com")
			So(err, ShouldBeNil)
			So(resp.Header.Get("Location"), ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("the status code, with Location headers for 201, 202", func() {
			codes := []int{201, 202}

			for _, code := range codes {
				maker.PlanGet("abc.com", code, "")
				resp, err := maker.Client().Get("abc.com")
				So(err, ShouldBeNil)
				So(resp.Header.Get("Location"), ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, code)
			}
		})
	})

	Convey("Repetition of Fakepoint default to a repetition of 1", t, func() {
		maker := NewFakepointMaker()

		Convey("defaults to 1", func() {
			maker.PlanGet("http://abc.com", 200, "")
			resp, err := maker.Client().Get("http://abc.com")
			resp2, err := maker.Client().Get("http://abc.com")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			So(resp2.StatusCode, ShouldEqual, 404)
		})

		Convey("is configurable", func() {
			maker.PlanGet("http://abc.com", 200, "").Duplicate(1)
			resp, err := maker.Client().Get("http://abc.com")
			resp2, err := maker.Client().Get("http://abc.com")
			resp3, err := maker.Client().Get("http://abc.com")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			So(resp2.StatusCode, ShouldEqual, 200)
			So(resp3.StatusCode, ShouldEqual, 404)
		})
	})

	Convey("It works with a Request object", t, func() {
		maker := NewFakepointMaker()
		maker.PlanGet("http://example.com", 200, "")
		req, err := http.NewRequest("GET", "http://example.com", nil)
		resp, err := maker.Client().Do(req)
		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, 200)
	})
}
