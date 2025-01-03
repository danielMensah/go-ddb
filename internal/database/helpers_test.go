package database

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestHelper_ParseCondition(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		args           []interface{}
		expectedKeys   map[string]types.AttributeValue
		expectedExpr   string
		expectedValues map[string]types.AttributeValue
		expectedErr    error
	}{
		{
			name:  "valid single condition",
			query: "userId = $1",
			args:  []interface{}{"123"},
			expectedKeys: map[string]types.AttributeValue{
				"userId": &types.AttributeValueMemberS{Value: "123"},
			},
			expectedExpr: "userId = :v1",
			expectedValues: map[string]types.AttributeValue{
				":v1": &types.AttributeValueMemberS{Value: "123"},
			},
			expectedErr: nil,
		},
		{
			name:  "valid multiple conditions",
			query: "userId = $1, age = $2",
			args:  []interface{}{"123", 25},
			expectedKeys: map[string]types.AttributeValue{
				"userId": &types.AttributeValueMemberS{Value: "123"},
				"age":    &types.AttributeValueMemberN{Value: "25"},
			},
			expectedExpr: "userId = :v1 AND age = :v2",
			expectedValues: map[string]types.AttributeValue{
				":v1": &types.AttributeValueMemberS{Value: "123"},
				":v2": &types.AttributeValueMemberN{Value: "25"},
			},
			expectedErr: nil,
		},
		{
			name:        "no arguments provided",
			query:       "userId = $1",
			args:        []interface{}{},
			expectedErr: ErrNoArguments,
		},
		{
			name:        "mismatched placeholders",
			query:       "userId = $1, age = $2",
			args:        []interface{}{"123"},
			expectedErr: ErrNoPlaceholders,
		},
		{
			name:        "invalid condition format",
			query:       "userId, age = $1",
			args:        []interface{}{"123", "222"},
			expectedErr: ErrConditionFormat,
		},
		{
			name:        "invalid placeholder format",
			query:       "userId = $x",
			args:        []interface{}{"123"},
			expectedErr: ErrPlaceholderFormat,
		},
		{
			name:        "unsupported argument type",
			query:       "userId = $1",
			args:        []interface{}{[]string{"123"}},
			expectedErr: ErrUnsupportedType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, expr, values, err := parseCondition(tt.query, tt.args...)

			// Assert the error
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Assert the keys, expression, and values
			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedKeys, keys)
				assert.Equal(t, tt.expectedExpr, expr)
				assert.Equal(t, tt.expectedValues, values)
			}
		})
	}
}
