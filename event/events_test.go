package event_test

import (
	. "github.com/reevoo/tracker/event"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/reevoo/tracker/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("NewEvent", func() {

	It("Generates an ID", func() {
		event := New(map[string][]string{})
		Expect(event["id"]).NotTo(BeNil())
	})

})