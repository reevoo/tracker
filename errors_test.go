package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var lastWrite string

type TestWriter struct{}

func (writer TestWriter) Write(p []byte) (n int, err error) {
	lastWrite = string(p[:])
	return len(lastWrite), nil
}

type TestSentryClient struct{}

func (client TestSentryClient) CaptureError(err error, tags map[string]string) string {
	var jsonTags []byte
	if tags != nil {
		jsonTags, _ = json.Marshal(tags)
		lastWrite = fmt.Sprintf("%s %s", err.Error(), string(jsonTags))
	} else {
		lastWrite = err.Error()
	}

	return ""
}

func TestErrors(t *testing.T) {
	err := errors.New("Sample error")

	Convey("Converts Go error into TrackerError", t, func() {
		Convey("Error messages should be equal", func() {
			trackerError := NewTrackerErrorFromError(err, nil)
			So(trackerError.Error(), ShouldEqual, "Sample error")
		})

		Convey("Should convert into map", func() {
			trackerError := NewTrackerErrorFromError(err, map[string]string{"foo": "bar"})
			m := trackerError.ToMap()
			So(m["name"], ShouldEqual, "Sample error")
			So(m["foo"], ShouldEqual, "bar")
		})
	})

	Convey("When ConsoleLogger is used", t, func() {
		logger := ConsoleLogger{writer: TestWriter{}}

		Convey("Logs error to the IO writer", func() {
			trackerError := NewTrackerErrorFromError(err, nil)
			logger.LogError(trackerError)
			So(lastWrite, ShouldEqual, "[ERROR] Sample error")
			lastWrite = ""
		})
	})

	Convey("When SentryErrorLogger is used", t, func() {
		logger := SentryErrorLogger{client: TestSentryClient{}}

		Convey("Logs error to Sentry Client", func() {
			trackerError := NewTrackerErrorFromError(err, nil)
			logger.LogError(trackerError)
			So(lastWrite, ShouldEqual, "Sample error")
			lastWrite = ""
		})

		Convey("Logs context map to Sentry Client", func() {
			trackerError := NewTrackerErrorFromError(err, map[string]string{"foo": "bar"})
			logger.LogError(trackerError)
			So(lastWrite, ShouldEqual, "Sample error {\"foo\":\"bar\",\"name\":\"Sample error\"}")
			lastWrite = ""
		})
	})
}
