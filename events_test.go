package tracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
)

var _ = Describe("NewEvent", func() {

	It("Generates an ID", func() {
		event := NewEvent("EventName", nil)
		Expect(event.Id).NotTo(BeNil())
	})

})
