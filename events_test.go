package tracker_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
)

var _ = Describe("NewEvent", func() {

	It("Generates an ID", func() {
		event := NewEvent(map[string][]string{})
		Expect(event["id"]).NotTo(BeNil())
	})

})
