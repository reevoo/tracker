package tracker_test

import (
	"bytes"
        "net/url"
	"errors"
	. "github.com/reevoo/tracker"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/reevoo/tracker/event"
	"net/http/httptest"
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

// Test implementation of EventStore.
type TestEventLogger struct {
	ThrowError bool
}

// Flips EventStored.
func (store TestEventLogger) Log(e interface{}) error {
	if store.ThrowError {
		return errors.New("TestEventStoreTriggeredError")
	}

	EventStored = true
        LastEvent = e.(event.Event)
	return nil
}

func EventToParams(event event.Event) string {
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

var _ = Describe("Server", func() {

	var (
		server   Server
		response *httptest.ResponseRecorder
		errors   = TestErrorLogger{}
		logger   = TestEventLogger{}
	)

	BeforeEach(func() {
		server = NewServer(ServerParams{
			ErrorLogger: &errors,
			EventLogger: &logger,
		})
	})

	Describe("GET /status", func() {
		BeforeEach(func() {
			response = get(&server, "/status")
		})

		It("returns HTTP Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("returns a string", func() {
			Expect(response.Body.String()).NotTo(BeEmpty())
		})

	})

	Describe("GET /event", func() {

		var (
			response *httptest.ResponseRecorder
			event    event.Event
			url      string
		)

		BeforeEach(func() {
			event = map[string][]string{
				"name": []string{"EventName"},
			}

			url = "/event?" + EventToParams(event)
		})

		It("returns HTTP 200", func() {
			response = get(&server, url)
			Expect(response.Code).To(Equal(200))
		})

		It("sends a request to the Event Store when JSON is correct", func() {
			EventStored = false

			response = get(&server, url)

			Eventually(func() bool {
				return EventStored
			}).Should(BeTrue())
		})

		It("creates an event with a UUID", func() {
			LastEvent = nil

			response = get(&server, url)

			Eventually(func() interface{} {
				if LastEvent == nil {
					return nil
				}
				return LastEvent["id"]
			}).ShouldNot(BeNil())
		})

		It("ignores any given UUID", func() {
			LastEvent = nil

			url = url + "&id=ID"

			response = get(&server, url)

			Eventually(func() bool {
				if LastEvent == nil {
					return false
				}

				return LastEvent["id"][0] != "ID"
			}).Should(BeTrue())
		})

		It("returns HTTP 400 when no params are given", func() {
			response = get(&server, "/event")
			Expect(response.Code).To(Equal(400))
		})

		It("tracks an error when the Event Store request fails", func() {
			logger.ThrowError = true
			ErrorThrown = false

			response = get(&server, url)

			Eventually(func() bool {
				return ErrorThrown
			}).Should(BeTrue())
			logger.ThrowError = false
		})
	})

})
