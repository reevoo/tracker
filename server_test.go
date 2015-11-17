package tracker_test

import (
	"bytes"
	"errors"
	. "github.com/reevoo/tracker"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
	"github.com/reevoo/tracker/event"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// Testing flag to check if an error is thrown.
var ErrorThrown = false

// Implementation of ErrorLogger that flips ErrorThrown.
type TestErrorLogger struct{}

// Flips ErrorThrown.
func (errorLogger TestErrorLogger) LogError(err TrackerError) {
	ErrorThrown = true
}

// Testing flag to check if an Event is stored.
var (
	EventStored = false
	LastEvent   event.Event
)

// Test implementation of EventLogger.
type TestEventLogger struct {
	ThrowError bool
}

// Flips EventStored.
func (store TestEventLogger) Log(e interface{}) error {
	if store.ThrowError {
		return errors.New("TestEventLoggerTriggeredError")
	}

	EventStored = true
	LastEvent = e.(event.Event)
	return nil
}

// Creates a URL query from event data
func eventToParams(event event.Event) string {
	var buf bytes.Buffer
	for key, values := range event {
		for _, value := range values {
			buf.WriteString(url.QueryEscape(key))
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(value))
			buf.WriteByte('&')
		}
	}
	s := buf.String()
	return s[0 : len(s)-1]
}

// Performs a GET request.
func get(server *Server, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)
	return resp
}

func TestServer(t *testing.T) {
	var errors = TestErrorLogger{}
	var logger = TestEventLogger{}
	var response *httptest.ResponseRecorder
	var server = NewServer(ServerParams{
		ErrorLogger: &errors,
		EventLogger: &logger,
		Environment: "production",
	})

	var eventuallyPollingInterval = 10 * time.Millisecond

	Convey("GET /status", t, func() {
		Convey("returns HTTP 200", func() {
			response = get(&server, "/status")
			So(response.Code, ShouldEqual, 200)
		})
		Convey("returns a non-empty string", func() {
			response = get(&server, "/status")
			So(response.Body.String(), ShouldNotBeNil)
		})
	})

	// Asynchronous Tests are implemented based on the discussions on the issue:
	// https://github.com/smartystreets/goconvey/issues/156
	Convey("GET /event", t, func() {
		event := map[string][]string{
			"name": {"EventName"},
		}

		url := "/event?" + eventToParams(event)

		Convey("returns HTTP 200", func() {
			response = get(&server, url)
			So(response.Code, ShouldEqual, 200)
		})

		Convey("sends a request to Logger", func(c C) {
			func() {
				EventStored = false
				response = get(&server, url)
			}()

			select {
			case <-time.After(eventuallyPollingInterval):
				c.So(EventStored, ShouldBeTrue)
			}
		})

		Convey("creates an event with a UUID", func(c C) {
			func() {
				LastEvent = nil
				response = get(&server, url)
			}()

			select {
			case <-time.After(eventuallyPollingInterval):
				c.So(LastEvent["id"][0], ShouldNotBeBlank)
			}
		})

		Convey("returns HTTP 400 when no params are given", func() {
			response = get(&server, "/event")
			So(response.Code, ShouldEqual, 400)
		})

		Convey("ignores any given UUID", func(c C) {
			func() {
				LastEvent = nil
				url = url + "&id=ID"
				response = get(&server, url)
			}()

			select {
			case <-time.After(eventuallyPollingInterval):
				c.So(LastEvent["id"][0], ShouldNotBeBlank)
				c.So(LastEvent["id"][0], ShouldNotEqual, "ID")
			}
		})

		Convey("tracks an error when the Logger request fails", func(c C) {
			func() {
				ErrorThrown = false
				logger.ThrowError = true
				response = get(&server, url)
			}()

			select {
			case <-time.After(eventuallyPollingInterval):
				c.So(ErrorThrown, ShouldBeTrue)
				logger.ThrowError = false
			}
		})
	})
}
