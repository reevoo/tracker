package tracker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"os"
)

var (
	DynamoDBTableName = os.Getenv("DYNAMODB_TABLE_NAME")
)

// An EventStore is used to permanently store events
type EventStore interface {
	Store(event Event) error
}

// We only care about PutItem in this program.
// This interface covers this method only to make mocking easier.
type dynamoDBPutter interface {
	PutItem(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
}

var (
	dynamoSession  = session.New(&aws.Config{Region: aws.String("eu-west-1")})
	DynamoDBClient *dynamodb.DynamoDB
)

func init() {
	DynamoDBClient = dynamodb.New(dynamoSession)
}

// Stores Events in DynamoDB.
type DynamoDBEventStore struct {
	DynamoDBClient dynamoDBPutter
}

// Stores an Event in DynamoDB.
func (store DynamoDBEventStore) Store(event Event) error {

	input := &dynamodb.PutItemInput{
		TableName: aws.String(DynamoDBTableName),
		Item:      event.ToDynamoDBItem(),
	}

	store.DynamoDBClient.PutItem(input)

	return nil
}

// Converts an Item to a DynamoDB map for the Item attribute of a PutItem call.
func (event Event) ToDynamoDBItem() map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"Id":       &dynamodb.AttributeValue{B: event.Id[:]},
		"Name":     &dynamodb.AttributeValue{S: &event.Name},
		"Metadata": &dynamodb.AttributeValue{M: nil}, // TODO
	}
}
