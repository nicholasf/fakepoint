package fake_round_trip

import (
	. "github.com/smartystreets/goconvey/convey"
//	"io/ioutil"
	"testing"
	"fmt"
)

func TestFakeRoundTrip(t *testing.T) {

//	Convey("Basic DSL - 3 arg creator func", t, func() {
//		client := NewFakeClient("GET", "https://api.opsgenie.com/v1/json/alert", "hello world")
//		Convey("It generates a fake client", func() {
//			So(client, ShouldNotBeNil)
//		})
//
//		Convey("it returns the document with a 200", func() {
//			resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
//			text, _ := ioutil.ReadAll(resp.Body)
//			So(string(text), ShouldEqual, "hello world")
//			So(resp.StatusCode, ShouldEqual, 200)
//		})
//
//		Convey("it doesnt resolve to the wrong URL", func() {
//			resp, _ := client.Get("https://somethingelse.com")
//			So(resp.StatusCode, ShouldEqual, 404)
//		})
//	})

	Convey("Chaining response setters on the client", t, func() {
		client := NewFakeClient("GET", "https://api.opsgenie.com/v1/json/alert", "hello world")

//		Convey("the header", func() {
//			client.SetResponseHeader("Content-Type", "application/json")
//			resp, err := client.Get("https://api.opsgenie.com/v1/json/alert")
//			So(err, ShouldBeNil)
//			So(resp.Header.Get("Content-Type"), ShouldEqual, "application/json")
//		})
//
//		Convey("the status code, generically", func() {
//			client = client.SetStatusCode(500)
//			resp, err := client.Get("https://api.opsgenie.com/v1/json/alert")
//			So(err, ShouldBeNil)
//			So(resp.StatusCode, ShouldEqual, 500)
//		})

		Convey("the status code, with Location headers for 201, 202, 302", func() {
//			locationRequiredCodes := []int{ 201, 302, 202 }
			locationRequiredCodes := []int{ 302 }

			for _, code := range locationRequiredCodes {
				fmt.Println("code: ", code)
				client = client.SetStatusCode(code)
				resp, err := client.Get("https://api.opsgenie.com/v1/json/alert")
				So(err, ShouldBeNil)
				So(resp.Header.Get("Location"), ShouldNotBeNil)
				fmt.Println("scode: ", resp.StatusCode)
				So(resp.StatusCode, ShouldEqual, code)
				fmt.Println("done.")
			}
		})


		//		Convey("the status code", func() {
		//			client = client.SetStatusCode(302)
		//		})

//		Convey("it can accept a header status code", t, func() {
//			resp, _ := client.Get("https://api.opsgenie.com/v1/json/alert")
//			foo, _ := ioutil.ReadAll(resp.Body)
//			So(resp, ShouldNotBeNil)
//			So(string(foo), ShouldEqual, "hello world")
//		})
	})
}
