package database

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// parseCondition parses the SQL-like condition string and arguments into a DynamoDB query.
func parseCondition(query string, args ...interface{}) (map[string]types.AttributeValue, string, map[string]types.AttributeValue, error) {
	if len(args) == 0 {
		return nil, "", nil, ErrNoArguments
	}

	var (
		keyAttributes          = make(map[string]types.AttributeValue) // For GetItem
		expressionValues       = make(map[string]types.AttributeValue) // For Query
		keyConditionExpression strings.Builder                         // For Query
		parts                  = strings.Split(query, ",")
	)

	if len(parts) != len(args) {
		return nil, "", nil, ErrNoPlaceholders
	}

	for i, part := range parts {
		conditionParts := strings.Split(strings.TrimSpace(part), "=")
		if len(conditionParts) != 2 {
			return nil, "", nil, ErrConditionFormat
		}

		attrName := strings.TrimSpace(conditionParts[0])
		placeholder := strings.TrimSpace(conditionParts[1])

		if placeholder != fmt.Sprintf("$%d", i+1) {
			return nil, "", nil, ErrPlaceholderFormat
		}

		key := fmt.Sprintf(":v%d", i+1)

		// Handle the argument and map it to an AttributeValue
		var attrValue types.AttributeValue
		switch v := args[i].(type) {
		case string:
			attrValue = &types.AttributeValueMemberS{Value: v}
		case int, int64, int32:
			attrValue = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", v)}
		default:
			return nil, "", nil, fmt.Errorf("%w: %v", ErrUnsupportedType, reflect.TypeOf(v))
		}

		// Populate the maps for GetItem and Query
		keyAttributes[attrName] = attrValue
		expressionValues[key] = attrValue

		// Build the KeyConditionExpression for Query
		if i > 0 {
			keyConditionExpression.WriteString(" AND ")
		}
		keyConditionExpression.WriteString(fmt.Sprintf("%s = %s", attrName, key))
	}

	return keyAttributes, keyConditionExpression.String(), expressionValues, nil
}

// buildUpdateExpression generates the UpdateExpression and ExpressionAttributeValues for an update operation.
func buildUpdateExpression(updates map[string]interface{}) (string, map[string]types.AttributeValue, error) {
	updateParts := make([]string, 0)
	expressionValues := make(map[string]types.AttributeValue)

	for field, value := range updates {
		placeholder := ":" + field
		updateParts = append(updateParts, fmt.Sprintf("%s = %s", field, placeholder))

		attrValue, err := attributevalue.Marshal(value)
		if err != nil {
			return "", nil, fmt.Errorf("marshalling update value for %s: %w", field, err)
		}
		expressionValues[placeholder] = attrValue
	}

	return "SET " + strings.Join(updateParts, ", "), expressionValues, nil
}
