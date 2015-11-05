package tracker

import (
	"fmt"
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
		"Metadata": &dynamodb.AttributeValue{M: event.Metadata.ToDynamoDBItem()},
	}
}

func (metadata Metadata) ToDynamoDBItem() map[string]*dynamodb.AttributeValue {
	var item = make(map[string]*dynamodb.AttributeValue)

	for key, value := range metadata {
		attrVal := attributeValue(value)
		item[key] = &attrVal
	}

	return item
}

// Convert something to a dynamodb.AttributeValue
// FIXME: What fresh hell is this?!
func attributeValue(item interface{}) dynamodb.AttributeValue {
	var value dynamodb.AttributeValue

	switch item.(type) {

	case []byte:
		value = dynamodb.AttributeValue{B: item.([]byte)}

	case bool:
		asBool := item.(bool)
		value = dynamodb.AttributeValue{BOOL: &asBool}

	case []interface{}:
		var inner []*dynamodb.AttributeValue

		for _, value := range item.([]interface{}) {
			attrVal := attributeValue(value)
			inner = append(inner, &attrVal)
		}

		value = dynamodb.AttributeValue{L: inner}

	case string:
		str := item.(string)
		value = dynamodb.AttributeValue{S: &str}

	case int:
		numStr := fmt.Sprint(item.(int))
		value = dynamodb.AttributeValue{N: &numStr}

	case float32:
		numStr := fmt.Sprint(item.(float32))
		value = dynamodb.AttributeValue{N: &numStr}

	case float64:
		numStr := fmt.Sprint(item.(float64))
		value = dynamodb.AttributeValue{N: &numStr}

	case map[string]interface{}:
		var inner = make(map[string]*dynamodb.AttributeValue)

		for key, value := range item.(map[string]interface{}) {
			attrVal := attributeValue(value)
			inner[key] = &attrVal
		}

		value = dynamodb.AttributeValue{M: inner}

	// TODO: Why isn't this a map[string]interface{}?
	case Metadata:
		var inner = make(map[string]*dynamodb.AttributeValue)

		for key, value := range map[string]interface{}(item.(Metadata)) {
			attrVal := attributeValue(value)
			inner[key] = &attrVal
		}

		value = dynamodb.AttributeValue{M: inner}

	case nil:
		null := true
		value = dynamodb.AttributeValue{NULL: &null}

	// TODO: Not implemented:
	// - BS [][]byte (Binary Set)
	// - NS []*string (Number Set)
	// - SS []*string (String Set)

	default:
		panic(fmt.Sprintf("Cannot parse item: %s", item))
	}

	return value
}
