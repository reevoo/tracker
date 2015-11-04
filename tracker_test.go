package tracker_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tracker Suite")
}

// HELPERS

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
