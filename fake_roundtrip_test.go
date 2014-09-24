package fake_client

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

func TestFakeRoundTrip(t *testing.T) {

	Convey("Basic DSL - 3 arg creator func", t, func() {
		client := NewFakeClient("GET", "https://api.opsgenie.com/v1/json/alert", "hello world")
		Convey("It generates a fake client", func() {
			So(client, ShouldNotBeNil)
		})

		Convey("it returns the document with a 200", func() {
			resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
			text, _ := ioutil.ReadAll(resp.Body)
			So(string(text), ShouldEqual, "hello world")
			So(resp.StatusCode, ShouldEqual, 200)
		})

		Convey("it doesnt resolve to the wrong URL", func() {
			resp, _ := client.Get("https://somethingelse.com")
			So(resp.StatusCode, ShouldEqual, 404)
		})
	})

	//
	//	Convey("it can accept a header status code", t, func() {
	//		resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
	//		foo, _ := ioutil.ReadAll(resp.Body)
	//		So(resp, ShouldNotBeNil)
	//		So(string(foo), ShouldEqual, "hello world")
	//	})
}
