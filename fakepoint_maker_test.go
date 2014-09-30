package fakepoint

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"io/ioutil"
)

func TestFakepointMaker(t *testing.T) {

	Convey("Create a Fake Client", t, func() {
		maker := NewFakepointMaker()
		So(maker, ShouldNotBeNil)
	})

	Convey("Create a http.Client", t, func() {
		maker := NewFakepointMaker()
		maker.PlanGet("https://api.opsgenie.com/v1/json/alert", 200, "hello world")

		resp, err := maker.Client().Get("https://api.opsgenie.com/v1/json/alert")
		text, err := ioutil.ReadAll(resp.Body)
		So(err, ShouldBeNil)
		So(string(text), ShouldEqual, "hello world")
		So(resp.StatusCode, ShouldEqual, 200)
	})
}
