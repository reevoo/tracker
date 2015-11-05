package tracker_test

import (
	"errors"
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
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

	Describe("POST /event", func() {

		var (
			response  *httptest.ResponseRecorder
			event     Event
			eventJson string
		)

		BeforeEach(func() {
			event = Event{
				Name:     "EventName",
				Metadata: make(map[string]interface{}),
			}

			eventJson = event.ToJson()
		})

		It("returns HTTP 200", func() {
			response = post(&server, "/event", eventJson)
			Expect(response.Code).To(Equal(200))
		})

		It("sends a request to DynamoDB when JSON is correct", func() {
			EventStored = false

			response = post(&server, "/event", eventJson)

			Eventually(func() bool {
				return EventStored
			}).Should(BeTrue())
		})

		It("creates an event with a UUID", func() {
			EventStored = false

			response = post(&server, "/event", eventJson)

			Eventually(func() uuid.UUID {
				return LastEvent.Id
			}).ShouldNot(BeNil())
		})

		It("return HTTP 200 when the event does not have metadata", func() {
			response = post(&server, "/event", Event{Name: "EventName", Metadata: nil}.ToJson())
			Expect(response.Code).To(Equal(200))
		})

		It("returns HTTP 400 when the event is not JSON", func() {
			response = post(&server, "/event", "Definitely Not JSON!")
			Expect(response.Code).To(Equal(400))
		})

		It("returns HTTP 400 when the event does not have a name", func() {
			response = post(&server, "/event", Event{Name: "", Metadata: make(map[string]interface{})}.ToJson())
			Expect(response.Code).To(Equal(400))
		})

		It("tracks an error when the DynamoDB request fails", func() {
			store.ThrowError = true
			ErrorThrown = false

			response = post(&server, "/event", eventJson)

			Eventually(func() bool {
				return ErrorThrown
			}).Should(BeTrue())
			store.ThrowError = false
		})
	})

})
