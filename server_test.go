package tracker_test

import (
	. "github.com/reevoo/tracker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func get(server *gin.Engine, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

var _ = Describe("Server", func() {

	var (
		server  *gin.Engine
	)

	BeforeEach(func() {
		server = CreateServer()
	})

	Describe("GET /status", func() {

		It("returns HTTP Status 200", func() {
			Expect(get(server, "/status").Code).To(Equal(200))
		})

	})

})
