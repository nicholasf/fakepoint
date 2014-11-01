package fakepoint

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
	"net/http"
)

const foo = "hello world"

func TestFakeRoundTrip(t *testing.T) {
	Convey("Basic DSL - planning roundtrips", t, func() {
		maker := NewFakepointMaker()

		Convey("it returns the document with a 200", func() {
			maker.NewGet("https://api.opsgenie.com/v1/json/alert", 200).SetHeader("Content-Type", "text/plain").SetResponse(foo)

			resp, _ := maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
			text, _ := ioutil.ReadAll(resp.Body)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "text/plain")
			So(string(text), ShouldEqual, "hello world")
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("it doesnt resolve to the wrong URL", func() {
			resp, _ := maker.Client().Get("https://somethingelse.com")
			So(resp.StatusCode, ShouldEqual, 404)
		})

		Convey("the header", func() {
			maker.NewGet("https://api.opsgenie.com/v1/json/alert", 200).SetHeader("Content-Type", "application/json").SetResponse("{ \"code\": 200 }")
			resp, err := maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
			So(err, ShouldBeNil)
			So(resp.Header.Get("Content-Type"), ShouldEqual, "application/json")
		})

		Convey("it returns errors when asked", func() {
			Convey("it returns proper http error codes when used", func() {
				maker.NewBadRequest("https://api.opsgenie.com/v1/json/alert", "GET")
				resp, err := maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, 400)

				maker.NewForbidden("https://api.opsgenie.com/v1/json/alert", "GET")
				resp, err = maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, 403)

				maker.NewUnauthorized("https://api.opsgenie.com/v1/json/alert", "GET")
				resp, err = maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, 401)

				maker.NewNotFound("https://api.opsgenie.com/v1/json/alert", "GET")
				resp, err = maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, 404)

				maker.NewInternalError("https://api.opsgenie.com/v1/json/alert", "GET")
				resp, err = maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
				So(err, ShouldBeNil)
				So(resp.StatusCode, ShouldEqual, 500)
			})
		})
	})

	Convey("Specify the response document", t, func() {
		maker := NewFakepointMaker()
		maker.NewGet("http://abc.com", 200).SetResponseDocument("./README.md")
		resp, err := maker.Client().Get("http://abc.com")
		text, _ := ioutil.ReadAll(resp.Body)
		So(err, ShouldBeNil)
		So(string(text), ShouldContainSubstring, "Create Fake endpoints for HTTP testing. Specify the response data sent back.")
	})

	Convey("The maker handles multiple fake round trips", t, func() {
		maker := NewFakepointMaker()
		Convey("distinct URLs", func() {
			maker.NewGet("https://api.opsgenie.com/v1/json/alert", 200).SetResponse("hello world")
			maker.NewGet("https://another.system.com", 200).SetResponse("not with a bang but a whimper")
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
			maker.NewPost("abc.com/greeting", 201).SetResponse("hello world1")
			maker.NewGet("abc.com/greeting", 200).SetResponse("hello world2")

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
			maker.NewGet("abc.com", 302).SetResponse("")
			resp, err := maker.Client().Get("abc.com")
			So(err, ShouldBeNil)
			So(resp.Header.Get("Location"), ShouldNotBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("the status code, with Location headers for 201, 202", func() {
			codes := []int{201, 202}

			for _, code := range codes {
				maker.NewGet("abc.com", code).SetResponse("")
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
			maker.NewGet("http://abc.com", 200).SetResponse("")
			resp, err := maker.Client().Get("http://abc.com")
			resp2, err := maker.Client().Get("http://abc.com")
			So(err, ShouldBeNil)
			So(resp.StatusCode, ShouldEqual, 200)
			So(resp2.StatusCode, ShouldEqual, 404)
		})

		Convey("is configurable", func() {
			maker.NewGet("http://abc.com", 200).SetResponse("").Duplicate(1)
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
		maker.NewGet("http://example.com", 200).SetResponse("hello world")
		req, err := http.NewRequest("GET", "http://example.com", nil)
		resp, err := maker.Client().Do(req)
		text, _ := ioutil.ReadAll(resp.Body)
		So(err, ShouldBeNil)
		So(string(text), ShouldEqual, "hello world")
		So(resp.StatusCode, ShouldEqual, 200)
	})
}
