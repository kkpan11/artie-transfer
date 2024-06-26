package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuoteLiteral(t *testing.T) {
	testCases := []struct {
		name     string
		colVal   string
		expected string
	}{
		{
			name:     "string",
			colVal:   "hello",
			expected: "'hello'",
		},
		{
			name:     "string that requires escaping",
			colVal:   "bobby o'reilly",
			expected: `'bobby o\'reilly'`,
		},
		{
			name:     "string with line breaks",
			colVal:   "line1 \n line 2",
			expected: "'line1 \n line 2'",
		},
		{
			name:     "string with existing backslash",
			colVal:   `hello \ there \ hh`,
			expected: `'hello \\ there \\ hh'`,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, QuoteLiteral(testCase.colVal), testCase.name)
	}
}

func TestQuoteLiterals(t *testing.T) {
	assert.Empty(t, QuoteLiterals(nil))
	assert.Empty(t, QuoteLiterals([]string{}))
	assert.Equal(t, []string{"'a'"}, QuoteLiterals([]string{"a"}))
	assert.Equal(t, []string{"'a'", "'b'", "'c\\'c'"}, QuoteLiterals([]string{"a", "b", "c'c"}))
}

func TestParseDataTypeDefinition(t *testing.T) {
	{
		dataType, parameters, err := ParseDataTypeDefinition("number")
		assert.NoError(t, err)
		assert.Equal(t, "number", dataType)
		assert.Empty(t, parameters)
	}
	{
		dataType, parameters, err := ParseDataTypeDefinition("number(5,2)")
		assert.NoError(t, err)
		assert.Equal(t, "number", dataType)
		assert.Equal(t, []string{"5", "2"}, parameters)
	}
	{
		dataType, parameters, err := ParseDataTypeDefinition("number(5, 2)")
		assert.NoError(t, err)
		assert.Equal(t, "number", dataType)
		assert.Equal(t, []string{"5", "2"}, parameters)
	}
	{
		dataType, parameters, err := ParseDataTypeDefinition("VARCHAR(1234)")
		assert.NoError(t, err)
		assert.Equal(t, "VARCHAR", dataType)
		assert.Equal(t, []string{"1234"}, parameters)
	}
	{
		// Spaces:
		dataType, parameters, err := ParseDataTypeDefinition("VARCHAR")
		assert.NoError(t, err)
		assert.Equal(t, "VARCHAR", dataType)
		assert.Empty(t, parameters)
	}
	{
		// Spaces + parameters:
		dataType, parameters, err := ParseDataTypeDefinition("   VARCHAR   (   1234   )  ")
		assert.NoError(t, err)
		assert.Equal(t, "VARCHAR", dataType)
		assert.Equal(t, []string{"1234"}, parameters)
	}
	{
		// Malformed parameters:
		_, _, err := ParseDataTypeDefinition("VARCHAR(1234")
		assert.ErrorContains(t, err, "missing closing parenthesis")
	}
}
