package event_test

import (
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/reevoo/tracker/event"
	"testing"
)

func TestEvents(t *testing.T) {

	Convey("Generating an ID", t, func() {
		event := New(map[string][]string{})
		So(event["id"], ShouldNotBeNil)
	})

	Convey("An Empty Event", t, func() {
		event := New(map[string][]string{})
		So(event.Empty(), ShouldBeTrue)
	})

}
