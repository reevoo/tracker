package tracker_test

import (
	"errors"
	. "github.com/reevoo/tracker"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/gomega"
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
	LastEvent   Event
)

// Test implementation of EventStore.
type TestEventStore struct {
	ThrowError bool
}

// Flips EventStored.
func (store TestEventStore) Store(event Event) error {
	if store.ThrowError {
		return errors.New("TestEventStoreTriggeredError")
	}

	EventStored = true
	LastEvent = event
	return nil
}

var _ = Describe("Server", func() {

	var (
		server   Server
		response *httptest.ResponseRecorder
		errors   = TestErrorLogger{}
		store    = TestEventStore{}
	)

	BeforeEach(func() {
		server = NewSilentServer(ServerParams{
			ErrorLogger: &errors,
			EventStore:  &store,
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
			event    Event
			url      string
		)

		BeforeEach(func() {
			event = map[string][]string{
				"name": []string{"EventName"},
			}

			url = "/event?" + event.ToParams()
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
				return LastEvent.Id()
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

				return LastEvent.Id() != "ID"
			}).Should(BeTrue())
		})

		It("returns HTTP 400 when no params are given", func() {
			response = get(&server, "/event")
			Expect(response.Code).To(Equal(400))
		})

		It("tracks an error when the Event Store request fails", func() {
			store.ThrowError = true
			ErrorThrown = false

			response = get(&server, url)

			Eventually(func() bool {
				return ErrorThrown
			}).Should(BeTrue())
			store.ThrowError = false
		})
	})

})
