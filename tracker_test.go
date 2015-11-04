package tracker_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Performs a GET request.
func get(server *Server, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

// Performs a POST request.
func post(server *Server, url string, body string) *httptest.ResponseRecorder {
	bodyReader := bytes.NewBufferString(body)
	req, _ := http.NewRequest("POST", url, bodyReader)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

// Hook the Go testing framework into Ginkgo and Gomega.
func TestTracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tracker Suite")
}
