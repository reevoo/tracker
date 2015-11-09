package tracker_test

import (
	. "github.com/reevoo/tracker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var LastWrite string

type TestStore struct{}

func (store TestStore) Write(p []byte) (n int, err error) {
	LastWrite = string(p[:])
	return len(LastWrite), nil
}

var _ = Describe("NewEventLog", func() {

	It("Writes to STDOUT by default", func() {
		eventLog := NewEventLog(nil)
		Expect(eventLog.Writer).To(Equal(os.Stdout))
	})

	Describe("Store()", func() {

		It("Writes JSON to writer", func() {
			testStore := TestStore{}
			event := NewEvent(map[string][]string{})
			eventLog := NewEventLog(testStore)

			eventLog.Store(event)
			Expect(LastWrite).To(Equal(event.ToJson()))
		})

	})
})
