package ecobank

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseError_Add(t *testing.T) {
	var err ResponseError
	err.Add("Error 1")
	err.Add("Error 2")

	assert.Equal(t, ResponseError([]string{"Error 1", "Error 2"}), err)
	assert.Equal(t, 2, err.Len())
}

func TestResponseError_Error(t *testing.T) {
	testCases := []struct {
		name     string
		errors   ResponseError
		expected string
	}{
		{
			name:     "no errors",
			errors:   ResponseError{},
			expected: "",
		},
		{
			name:     "single error",
			errors:   ResponseError([]string{"Single Error"}),
			expected: "Single Error",
		},
		{
			name:     "multiple errors",
			errors:   ResponseError([]string{"Error A", "Error B", "Error C"}),
			expected: "Error A\nError B\nError C",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.errors.Error())
		})
	}
}

func TestResponseError_GetErrors(t *testing.T) {
	errors := ResponseError([]string{"Error 1", "Error 2"})
	result := errors.All()

	assert.Equal(t, []string{"Error 1", "Error 2"}, result)
	assert.Equal(t, errors.Len(), len(result))
}

func TestResponseError_String(t *testing.T) {
	testCases := []struct {
		name     string
		errors   ResponseError
		expected string
	}{
		{
			name:     "no errors",
			errors:   ResponseError{},
			expected: "",
		},
		{
			name:     "single error",
			errors:   ResponseError([]string{"Single Error"}),
			expected: "Single Error",
		},
		{
			name:     "multiple errors",
			errors:   ResponseError([]string{"Error A", "Error B", "Error C"}),
			expected: "Error A\nError B\nError C",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.errors.String())
		})
	}
}

func TestResponseError_Empty(t *testing.T) {
	var err ResponseError
	assert.Equal(t, "", err.Error())
	assert.Nil(t, err.All())
	assert.Equal(t, "", err.String())
}

func TestResponseError_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		errors   ResponseError
		expected string
	}{
		{
			name:     "no errors",
			errors:   ResponseError{},
			expected: "[]",
		},
		{
			name:     "single error",
			errors:   ResponseError([]string{"Single Error"}),
			expected: `["Single Error"]`,
		},
		{
			name:     "multiple errors",
			errors:   ResponseError([]string{"Error A", "Error B", "Error C"}),
			expected: `["Error A","Error B","Error C"]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marshaled, err := json.Marshal(tc.errors)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(marshaled))
		})
	}
}
