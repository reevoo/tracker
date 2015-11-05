package tracker_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/reevoo/tracker"
)

var LastInput dynamodb.PutItemInput

type TestDynamoDBClient struct {
	PutItemOutput dynamodb.PutItemOutput
}

func (client TestDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	LastInput = *input
	return &client.PutItemOutput, nil
}

var _ = Describe("dynamodb file init", func() {
	It("creates a client by default", func() {
		Expect(DynamoDBClient).NotTo(BeNil())
	})
})

var _ = Describe("Event", func() {

	Describe("ToDynamoDBItem", func() {

		var (
			event Event
			input map[string]*dynamodb.AttributeValue
		)

		BeforeEach(func() {
			event = NewEvent("TestEventName", nil)
			input = event.ToDynamoDBItem()
		})

		It("Parses the Id", func() {
			idFromInput, _ := uuid.Parse(input["Id"].B)
			Expect(*idFromInput).To(Equal(event.Id))
		})

		It("Parses the Name", func() {
			Expect(*input["Name"].S).To(Equal("TestEventName"))
		})

		It("Parses an empty Metadata", func() {
			event = NewEvent("TestEventName", nil)
			input = event.ToDynamoDBItem()

			Expect(input["Metadata"].M).To(Equal(map[string]*dynamodb.AttributeValue{}))
		})

		It("Parses a non-empty Metadata", func() {
			event = NewEvent("TestEventName", ExampleMetadata)
			input = event.ToDynamoDBItem()

			// TODO: More specific testing of data types.

			Expect(input["Metadata"].M).NotTo(BeNil())
		})

	})

})

var _ = Describe("DynamoDBEventStore", func() {

	var (
		store  DynamoDBEventStore
		client TestDynamoDBClient
	)

	BeforeEach(func() {
		DynamoDBTableName = "TestTableName"
		client = TestDynamoDBClient{}

		store = DynamoDBEventStore{
			DynamoDBClient: &client,
		}
	})

	Describe(".Store()", func() {

		It("Sends the input to the client", func() {
			event := NewEvent("TestEventName", nil)
			expectedInput := dynamodb.PutItemInput{
				TableName: aws.String("TestTableName"),
				Item:      event.ToDynamoDBItem(),
			}

			err := store.Store(event)

			Expect(err).To(BeNil())
			Expect(LastInput).To(Equal(expectedInput))
		})

	})

})
