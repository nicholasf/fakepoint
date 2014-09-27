package fake_roundtrip

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFakeRoundTripClient(t *testing.T) {

	Convey("Create a Fake Client, without a base domain", t, func() {
		Convey("without a domain", func() {
			client := NewFakeClient()
			So(client, ShouldNotBeNil)
		})
	})
}
