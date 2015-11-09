package tracker_test

import (
	. "github.com/reevoo/tracker"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("NewEvent", func() {

	It("Generates an ID", func() {
		event := NewEvent(map[string][]string{})
		Expect(event["id"]).NotTo(BeNil())
	})

})
