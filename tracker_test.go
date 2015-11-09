package tracker_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	// Disable verbose logging on Test.
	fmt.Print("NOTE: Verbose logging is disabled on test.\nTo change this, go to `tracker_test.go`.\n\n")
	SetServerMode("release")
}

// Performs a GET request.
func get(server *Server, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()
	server.ServeHTTP(resp, req)

	return resp
}

// Hook the Go testing framework into Ginkgo and Gomega.
func TestTracker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tracker Suite")
}
