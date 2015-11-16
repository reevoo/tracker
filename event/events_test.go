package event_test

import (
  "testing"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	. "github.com/reevoo/tracker/event"
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
