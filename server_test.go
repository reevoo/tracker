package tracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
	"net/http/httptest"
)

// Logs errors to an exposed array.
type TestErrorLogger struct {
	LastError TrackerError
}

// Logs an error to an exposed array.
func (errorLogger TestErrorLogger) LogError(err TrackerError) {
	errorLogger.LastError = err
}

var _ = Describe("Server", func() {

	var (
		server   Server
		response *httptest.ResponseRecorder
		errors   = TestErrorLogger{}
	)

	BeforeEach(func() {
		server = NewServer(errors)
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
			response         *httptest.ResponseRecorder
			validRequestJson string
		)

		BeforeEach(func() {
			validRequestJson = Event{Name: "EventName", Metadata: make(map[string]interface{})}.ToJson()
		})

		It("returns HTTP 200", func() {
			response = post(&server, "/event", validRequestJson)
			Expect(response.Code).To(Equal(200))
		})

		PIt("sends a request to DynamoDB when JSON is correct", func() {
			response = post(&server, "/event", validRequestJson)
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

		PIt("tracks an error when the DynamoDB request fails", func() {
			response = post(&server, "/event", validRequestJson)

			Eventually(func() TrackerError {
				return errors.LastError
			}).ShouldNot(BeNil())
		})
	})

})
