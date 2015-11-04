package tracker_test

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
	"net/http"
	"net/http/httptest"
)

func get(server *gin.Engine, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

func post(server *gin.Engine, url string, body string) *httptest.ResponseRecorder {
	bodyReader := bytes.NewBufferString(body)
	req, _ := http.NewRequest("POST", url, bodyReader)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

var _ = Describe("Server", func() {

	var (
		server *gin.Engine
	)

	BeforeEach(func() {
		server = NewTrackerEngine()
	})

	Describe("GET /status", func() {

		var (
			response *httptest.ResponseRecorder
		)

		BeforeEach(func() {
			response = get(server, "/status")
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
			response            *httptest.ResponseRecorder
			validRequestJson    string
			requestsPerSecond = 1000
		)

		BeforeEach(func() {
			validRequestJson = Event{Name: "EventName", Metadata: make(map[string]interface{})}.ToJson()
			Errors.Clear()
		})

		It("returns HTTP 200", func() {
			response = post(server, "/event", validRequestJson)
			Expect(response.Code).To(Equal(200))
		})

		It("sends a request to DynamoDB when JSON is correct", func() {
			response = post(server, "/event", validRequestJson)
		})

		It("return HTTP 200 when the event does not have metadata", func() {
			response = post(server, "/event", Event{Name: "EventName", Metadata: nil}.ToJson())
			Expect(response.Code).To(Equal(200))
		})

		It("returns HTTP 400 when the event is not JSON", func() {
			response = post(server, "/event", "Definitely Not JSON!")
			Expect(response.Code).To(Equal(400))
		})

		It("returns HTTP 400 when the event does not have a name", func() {
			response = post(server, "/event", Event{Name: "", Metadata: make(map[string]interface{})}.ToJson())
			Expect(response.Code).To(Equal(400))
		})

		It("tracks an error when the DynamoDB request fails", func() {
			response = post(server, "/event", validRequestJson)

			Eventually(func() int {
				return Errors.Count()
			}).Should(Equal(1))
		})

		It(fmt.Sprintf("can handle %d requests in 1 second", requestsPerSecond), func() {
			for i := 0; i < requestsPerSecond; i++ {
				go post(server, "/event", validRequestJson)
			}

			Eventually(func() int {
				return Errors.Count()
			}).Should(Equal(requestsPerSecond))
		})
	})

})
