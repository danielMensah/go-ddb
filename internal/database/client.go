package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	ErrQuery               = errors.New("error querying DynamoDB")
	ErrGetItem             = errors.New("error getting item from DynamoDB")
	ErrNoArguments         = errors.New("no arguments provided for query")
	ErrNoPlaceholders      = errors.New("number of placeholders does not match number of arguments")
	ErrConditionFormat     = errors.New("each condition must be in the format 'attribute = $index'")
	ErrPlaceholderFormat   = errors.New("placeholder must be in the format '$index'")
	ErrUnsupportedType     = errors.New("unsupported argument type")
	ErrGetItemKeyCondition = errors.New("GetItem requires exactly one key condition")
	ErrParameterNotPointer = errors.New("out parameter must be a pointer")
)

// Client is an interface for the DynamoDB client.
type Client struct {
	db        DynamoDBClient
	tableName *string
}

// NewClient returns a new DynamoDB client.
func NewClient(db DynamoDBClient, tableName *string) *Client {
	return &Client{db, tableName}
}

// CreateItem handles creating a new item in the database.
func (c *Client) CreateItem(ctx context.Context, in interface{}) (*dynamodb.PutItemOutput, error) {
	item, err := attributevalue.MarshalMap(in)
	if err != nil {
		return nil, fmt.Errorf("marshalling input: %w", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: c.tableName,
	}

	return c.db.PutItem(ctx, input)
}

// Find handles querying the database for items.
func (c *Client) Find(ctx context.Context, out interface{}, query string, args ...interface{}) error {
	outType := reflect.TypeOf(out)
	if outType.Kind() != reflect.Ptr {
		return ErrParameterNotPointer
	}
	outElem := outType.Elem()
	if outElem.Kind() == reflect.Slice {
		return c.queryItems(ctx, out, query, args...)
	}

	return c.queryItem(ctx, out, query, args...)
}

func (c *Client) queryItems(ctx context.Context, out interface{}, query string, args ...interface{}) error {
	_, keyConditionExpression, expressionAttributeValues, err := parseCondition(query, args...)
	if err != nil {
		return err
	}

	input := &dynamodb.QueryInput{
		TableName:                 c.tableName,
		KeyConditionExpression:    &keyConditionExpression,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	result, err := c.db.Query(ctx, input)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrQuery, err)
	}

	return attributevalue.UnmarshalListOfMaps(result.Items, out)
}

func (c *Client) queryItem(ctx context.Context, out interface{}, query string, args ...interface{}) error {
	attributeValues, _, _, err := parseCondition(query, args...)
	if err != nil {
		return err
	}

	if len(attributeValues) != 1 {
		return ErrGetItemKeyCondition
	}

	input := &dynamodb.GetItemInput{
		TableName: c.tableName,
		Key:       attributeValues,
	}

	result, err := c.db.GetItem(ctx, input)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrGetItem, err)
	}

	return attributevalue.UnmarshalMap(result.Item, out)
}

// UpdateItem updates an existing item in the database.
func (c *Client) UpdateItem(ctx context.Context, key interface{}, updates map[string]interface{}) error {
	keyMap, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("marshalling key: %w", err)
	}

	updateExpression, expressionValues, err := buildUpdateExpression(updates)
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 c.tableName,
		Key:                       keyMap,
		UpdateExpression:          &updateExpression,
		ExpressionAttributeValues: expressionValues,
	}

	_, err = c.db.UpdateItem(ctx, input)
	if err != nil {
		return fmt.Errorf("updating item: %w", err)
	}

	return nil
}

// DeleteItem deletes an item from the database.
func (c *Client) DeleteItem(ctx context.Context, key interface{}) error {
	keyMap, err := attributevalue.MarshalMap(key)
	if err != nil {
		return fmt.Errorf("marshalling key: %w", err)
	}

	input := &dynamodb.DeleteItemInput{
		TableName: c.tableName,
		Key:       keyMap,
	}

	_, err = c.db.DeleteItem(ctx, input)
	if err != nil {
		return fmt.Errorf("deleting item: %w", err)
	}

	return nil
}
